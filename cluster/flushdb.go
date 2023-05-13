package cluster

import (
	"GoRedis/interface/resp"
	"GoRedis/resp/reply"
)

func flushdb(d *Database, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	replies := d.broadcast(c, cmdArgs)
	var errReply reply.ErrorReply
	for _, r := range replies {
		if reply.IsErrReply(r) {
			errReply = r.(reply.ErrorReply)
			break
		}
	}
	if errReply == nil {
		return reply.MakeOkReply()
	}
	return reply.MakeErrReply("Err " + errReply.Error())
}
