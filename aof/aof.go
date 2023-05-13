package aof

import (
	"GoRedis/config"
	databaseface "GoRedis/interface/database"
	"GoRedis/lib/logger"
	"GoRedis/lib/utils"
	"GoRedis/resp/parser"
	"GoRedis/resp/reply"
	"github.com/hdt3213/godis/redis/connection"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

// CmdLine is alias for [][]byte, represents a command line
type CmdLine = [][]byte

const (
	aofQueueSize = 1 << 16
)

type payload struct {
	cmdLine CmdLine
	dbIndex int
}

// Handler receive msgs from channel and write to AOF file
type Handler struct {
	enable      bool
	database    databaseface.IDatabase
	aofChan     chan *payload
	aofFile     *os.File
	aofFilename string
	// aof goroutine will send msg to main goroutine through this channel when aof tasks finished and ready to shutdown
	aofFinished chan struct{}
	// pause aof for start/finish aof rewrite progress
	pausingAof sync.RWMutex
	currentDB  int
}

func NewAofHandler(database databaseface.IDatabase) (*Handler, error) {
	handler := &Handler{}
	handler.enable = config.Properties.AppendOnly
	if !handler.enable {
		return nil, nil
	}
	handler.aofFilename = config.Properties.AppendFilename
	handler.database = database

	file, err := os.OpenFile(handler.aofFilename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	handler.aofChan = make(chan *payload, aofQueueSize)
	handler.aofFile = file

	go handler.handleAof()
	handler.loadAof() // recover data

	return handler, nil
}

func (h *Handler) AddAof(dbIndex int, cmd CmdLine) {
	if h.enable {
		h.aofChan <- &payload{
			cmdLine: cmd,
			dbIndex: dbIndex,
		}
	}
}

// aof data from handler.aofChan
func (h *Handler) handleAof() {
	h.currentDB = 0
	var lastFailed []byte
	retryCount := 1
	for p := range h.aofChan {
		// handler last time exit unexpectedly
		if lastFailed != nil {
			_, err := h.aofFile.Write(lastFailed)
			if err == nil {
				lastFailed = nil
				retryCount = 1
			} else {
				logger.Error("retry count "+strconv.Itoa(retryCount)+"failed ,error :", err)
				retryCount++
				continue
			}
		}

		// select standalone
		if p.dbIndex != h.currentDB {
			data := reply.MakeMultiBulkReply(utils.ToCmdLine("select", strconv.Itoa(p.dbIndex))).ToBytes()
			_, err := h.aofFile.Write(data)
			if err != nil {
				logger.Fatal("aof persistence error :", err)
				lastFailed = data
				continue
			}
			h.currentDB = p.dbIndex
		}

		data := reply.MakeMultiBulkReply(p.cmdLine).ToBytes()
		_, err := h.aofFile.Write(data)
		if err != nil {
			logger.Fatal("aof persistence error :", err)
			lastFailed = data
			continue
		}
	}
}

// recover data
func (h *Handler) loadAof() {
	logger.Info("data recover start")
	file, err := os.Open(h.aofFilename)
	if err != nil {
		logger.Error(err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	ch := parser.ParseStream(file)
	fakeConn := &connection.FakeConn{} // only used for save dbIndex
	defer logger.Info("data recover completed")
	for {
		select {
		case p := <-ch:
			if p.Err != nil {
				if p.Err == io.EOF {
					break
				}
				logger.Error(p.Err)
				continue
			}

			if p.Data == nil {
				//logger.Error("empty payload")
				continue
			}

			r, ok := p.Data.(*reply.MultiBulkReply)
			if !ok {
				logger.Error("need multi bulk")
				continue
			}

			rep := h.database.Exec(fakeConn, r.Args)
			if reply.IsErrReply(rep) {
				logger.Error(rep)
			}
		case <-time.After(200 * time.Millisecond):
			return
		}
	}
}
