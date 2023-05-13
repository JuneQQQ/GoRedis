package handler

/*
 * A tcp.RespHandler implements redis protocol
 */

import (
	"GoRedis/cluster"
	"GoRedis/config"
	databaseface "GoRedis/interface/database"
	"GoRedis/interface/resp"
	"GoRedis/lib/logger"
	"GoRedis/resp/parser"
	"GoRedis/resp/reply"
	"GoRedis/standalone"
	"context"
	"github.com/hdt3213/godis/redis/connection"
	"io"
	"net"
	"strings"
	"sync"
	atomic2 "sync/atomic"
)

var (
	unknownErrReplyBytes = []byte("-Err unknown\r\n")
)

// RespHandler implements tcp.Handler and serves as a redis handler
type RespHandler struct {
	activeConn sync.Map // *client -> placeholder
	db         databaseface.IDatabase
	closing    atomic2.Bool // refusing new client and new request
}

// MakeHandler creates a RespHandler instance
func MakeHandler() *RespHandler {
	var db databaseface.IDatabase
	//db = standalone.NewStandaloneDatabase()
	if config.Properties.Self != "" && len(config.Properties.Peers) > 0 {
		logger.Info("staring through cluster")
		db = cluster.MakeClusterDatabase()
	} else {
		logger.Info("staring through standalone")
		db = standalone.NewStandaloneDatabase()
	}
	return &RespHandler{
		db: db,
	}
}

func (h *RespHandler) closeClient(client *connection.Connection) {
	_ = client.Close()
	h.db.AfterClientClose(client)
	h.activeConn.Delete(client)
}

// Handle receives and executes redis commands
func (h *RespHandler) Handle(ctx context.Context, c net.Conn) {
	if h.closing.Load() {
		// closing handler refuse new connection
		_ = c.Close()
	}

	// abstract Redis connection from network connection
	conn := connection.NewConn(c)
	h.activeConn.Store(conn, 1)

	ch := parser.ParseStream(c)
	for payload := range ch {
		// handle error
		if payload.Err != nil {
			if payload.Err == io.EOF ||
				payload.Err == io.ErrUnexpectedEOF ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {
				// connection closed
				h.closeClient(conn)
				logger.Info("connection closed: " + conn.RemoteAddr().String())
				return
			}
			// protocol err
			errReply := reply.MakeErrReply(payload.Err.Error())
			err := conn.Write(errReply.ToBytes())
			if err != nil {
				h.closeClient(conn)
				logger.Info("connection closed: " + conn.RemoteAddr().String())
				return
			}
			continue
		}
		// assert not empty
		if payload.Data == nil {
			//logger.Error("empty payload")
			continue
		}
		var result resp.Reply
		switch payload.Data.(type) {
		case *reply.MultiBulkReply:
			// correctly parsed command, exec it by handler
			result = h.db.Exec(conn, payload.Data.(*reply.MultiBulkReply).Args)
		case *reply.ProtocolErrReply:
			result = payload.Data.(*reply.ProtocolErrReply)
		case *reply.SyntaxErrReply:
			result = payload.Data.(*reply.SyntaxErrReply)
		default:
			result = reply.MakeUnknownErrReply()
		}
		_ = conn.Write(result.ToBytes())

		// prepare handle incoming data
		//r, ok := payload.Data.(*reply.MultiBulkReply)
		//if !ok {
		//	logger.Error("require multi bulk reply")
		//	continue
		//}

	}
}

// Close stops handler
func (h *RespHandler) Close() error {
	logger.Info("handler shutting down...")
	h.closing.Store(true)
	// TODO: concurrent wait
	h.activeConn.Range(func(key interface{}, val interface{}) bool {
		client := key.(*connection.Connection)
		_ = client.Close()
		return true
	})
	h.db.Close()
	return nil
}
