package cluster

import (
	"GoRedis/interface/resp"
	"GoRedis/resp/reply"
)

func rename(d *Database, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	if len(cmdArgs) != 3 {
		return reply.MakeErrReply("wrong number of args")
	}

	src := string(cmdArgs[1])
	dst := string(cmdArgs[2])

	n1 := d.peerPicker.GetNode(src)
	n2 := d.peerPicker.GetNode(dst)

	if n1 != n2 {
		return reply.MakeErrReply("Err rename must in the same node")
	}
	return d.relay(n1, c, cmdArgs)
}
