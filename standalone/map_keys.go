package standalone

import (
	"GoRedis/interface/resp"
	"GoRedis/lib/utils"
	"GoRedis/resp/reply"
	"github.com/hdt3213/godis/lib/wildcard"
)

// Remove
func execDel(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}
	deleted := db.RemoveBulk(keys...)

	db.addAof(utils.ToCmdLine2("del", args...))
	return reply.MakeIntReply(int64(deleted))
}

func execExists(db *DB, args [][]byte) resp.Reply {
	result := int64(0)
	for _, arg := range args {
		key := string(arg)
		_, exists := db.GetEntity(key)
		if exists {
			result++
		}
	}

	return reply.MakeIntReply(result)
}

func execFlushDB(db *DB, args [][]byte) resp.Reply {
	db.data.Clear()

	db.addAof(utils.ToCmdLine2("flushdb", args...))
	return reply.MakeOkReply()
}

func execType(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.MakeStatusReply("key not exists")
	}
	switch entity.Data.(type) {
	case []byte:
		return reply.MakeStatusReply("string")
		// TODO other data structure
	}
	return reply.UnknownErrReply{}
}

// may be overwritten old entity
func execRename(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dst := string(args[1])
	entity, exists := db.GetEntity(src)
	if !exists {
		return reply.MakeErrReply("no such key")
	}
	db.Remove(dst)
	db.SetEntity(dst, entity)
	db.Remove(src)

	db.addAof(utils.ToCmdLine2("rename", args...))
	return reply.MakeIntReply(1)
}

func execRenameNx(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dst := string(args[1])
	_, exists := db.GetEntity(dst)
	if !exists {
		return reply.MakeIntReply(0)
	}
	entities, exists := db.GetEntity(src)
	db.SetEntity(dst, entities)
	db.Remove(src)

	db.addAof(utils.ToCmdLine2("renamenx", args...))
	return reply.MakeIntReply(1)
}

// keys *
func execKeys(db *DB, args [][]byte) resp.Reply {
	pattern := wildcard.CompilePattern(string(args[0]))
	result := make([][]byte, 0)
	db.data.ForEach(func(key string, val any) bool {
		if pattern.IsMatch(key) {
			result = append(result, []byte(key))
		}
		return true
	})
	return reply.MakeMultiBulkReply(result)
}

func init() {
	RegisterCommand("del", execDel, -2)
	RegisterCommand("exists", execExists, -2)
	RegisterCommand("flushdb", execFlushDB, 1)
	RegisterCommand("type", execType, 2)
	RegisterCommand("rename", execRename, 3)
	RegisterCommand("renamenx", execRenameNx, 3)
	RegisterCommand("keys", execKeys, 2)
}
