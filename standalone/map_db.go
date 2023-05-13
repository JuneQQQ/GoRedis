package standalone

import (
	"GoRedis/datastructure/dict"
	"GoRedis/interface/database"
	"GoRedis/interface/resp"
	"GoRedis/resp/reply"
	"strings"
)

type DB struct {
	index  int
	data   dict.Dict
	addAof func(line CmdLine)
}

func (db *DB) Exec(client resp.Connection, args database.CmdLine) resp.Reply {
	cmdName := strings.ToLower(string(args[0]))
	cmd := cmdTable[cmdName]
	if cmd == nil {
		return reply.MakeErrReply("Err unknown command " + cmdName)
	}
	// set k v
	if !validateArity(cmd.arity, args) {
		return reply.MakeArgNumErrReply(cmdName)
	}
	executor := cmd.executor
	return executor(db, args[1:])
}

func (db *DB) AfterClientClose(c resp.Connection) {
	//TODO implement me
	panic("implement me")
}

func (db *DB) Close() {
	//TODO implement me
	panic("implement me")
}

type ExecFunc func(db *DB, args [][]byte) resp.Reply
type CmdLine [][]byte

func makeDB() *DB {
	return &DB{data: dict.MakeSyncDict(), addAof: func(line CmdLine) {
		// do nothing to prevent inserting the same data when redis starting
	}}
}

func validateArity(arity int, cmdArgs [][]byte) bool {
	argNum := len(cmdArgs)
	if arity >= 0 {
		return argNum == arity
	}
	return argNum >= -arity
}

/* ---- data access ----- */

// GetEntity returns DataEntity bind to given key
func (db *DB) GetEntity(key string) (*database.DataEntity, bool) {
	raw, ok := db.data.Get(key)
	if !ok {
		return nil, false
	}
	entity, _ := raw.(*database.DataEntity)
	return entity, true
}

// SetEntity a DataEntity into DB
func (db *DB) SetEntity(key string, entity *database.DataEntity) int {
	return db.data.Set(key, entity)
}

// SetIfExists edit an existing DataEntity
func (db *DB) SetIfExists(key string, entity *database.DataEntity) int {
	return db.data.SetIfExists(key, entity)
}

// SetIfAbsent insert an DataEntity only if the key not exists
func (db *DB) SetIfAbsent(key string, entity *database.DataEntity) int {
	return db.data.SetIfAbsent(key, entity)
}

// Remove the given key from db
func (db *DB) Remove(key string) {
	db.data.Remove(key)
}

// RemoveBulk the given keys from db
func (db *DB) RemoveBulk(keys ...string) (deleted int) {
	deleted = 0
	for _, key := range keys {
		_, exists := db.data.Get(key)
		if exists {
			db.Remove(key)
			deleted++
		}
	}
	return deleted
}

// Flush clean standalone
func (db *DB) Flush() {
	db.data.Clear()
}
