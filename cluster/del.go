package cluster

import (
	"GoRedis/interface/resp"
	"GoRedis/resp/reply"
)

func del(d *Database, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	replies := d.broadcast(c, cmdArgs)
	var deleted int64
	var errReply reply.ErrorReply
	for _, r := range replies {
		intReply, ok := r.(*reply.IntReply)
		if reply.IsErrReply(r) || !ok {
			errReply = r.(reply.ErrorReply)
			break
		}
		deleted += intReply.Code

	}
	if errReply == nil {
		return reply.MakeIntReply(deleted)
	}
	return reply.MakeErrReply("Err " + errReply.Error())
}
