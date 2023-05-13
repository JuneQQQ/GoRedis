package cluster

import "GoRedis/interface/resp"

func ping(d *Database, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	return d.db.Exec(c, cmdArgs)
}
