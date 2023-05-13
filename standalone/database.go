package standalone

import (
	"GoRedis/aof"
	"GoRedis/config"
	"GoRedis/interface/resp"
	"GoRedis/lib/logger"
	"GoRedis/resp/reply"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
)

type Database struct {
	dbSet      []*DB
	aofHandler *aof.Handler
}

// NewStandaloneDatabase creates a redis standalone,
func NewStandaloneDatabase() *Database {
	database := &Database{}
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}

	database.dbSet = make([]*DB, config.Properties.Databases)
	for i := range database.dbSet {
		singleDB := makeDB()
		singleDB.index = i
		database.dbSet[i] = singleDB
	}

	if config.Properties.AppendOnly {
		// 1. load aof handler
		// 2. recover old data
		// must execute before aofFunc initialize
		handler, err := aof.NewAofHandler(database)
		if err != nil {
			panic(err)
		}
		database.aofHandler = handler
	}

	// init aofFunc
	for _, set := range database.dbSet {
		set.addAof = func(line CmdLine) {
			sdb := set
			// Pay attention to closure issues
			database.aofHandler.AddAof(sdb.index, line)
		}
	}

	return database
}

// Exec executes command
// parameter `cmdLine` contains command and its arguments, for example: "set key value"
func (mdb *Database) Exec(c resp.Connection, cmdLine [][]byte) (result resp.Reply) {
	defer func() {
		if err := recover(); err != nil {
			logger.Warn(fmt.Sprintf("error occurs: %v\n%s", err, string(debug.Stack())))
		}
	}()

	cmdName := strings.ToLower(string(cmdLine[0]))
	if cmdName == "select" {
		if len(cmdLine) != 2 {
			return reply.MakeArgNumErrReply("select")
		}
		return execSelect(c, mdb, cmdLine[1:])
	}
	// normal commands
	dbIndex := c.GetDBIndex()
	selectedDB := mdb.dbSet[dbIndex]
	return selectedDB.Exec(c, cmdLine)
}

// Close graceful shutdown standalone
func (mdb *Database) Close() {

}

func (mdb *Database) AfterClientClose(c resp.Connection) {
}

func execSelect(c resp.Connection, mdb *Database, args [][]byte) resp.Reply {
	dbIndex, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.MakeErrReply("Err invalid DB index")
	}
	if dbIndex >= len(mdb.dbSet) {
		return reply.MakeErrReply("Err DB index is out of range")
	}
	c.SelectDB(dbIndex)
	return reply.MakeOkReply()
}
