package database

import (
	"GoRedis/interface/resp"
)

// CmdLine is alias for [][]byte, represents a command line
type CmdLine = [][]byte

// IDatabase is the interface for redis style storage engine
type IDatabase interface {
	Exec(conn resp.Connection, args CmdLine) resp.Reply
	AfterClientClose(c resp.Connection)
	Close()
}

// DataEntity stores data bound to a key, including a string, list, hash, set and so on
type DataEntity struct {
	Data any
}
