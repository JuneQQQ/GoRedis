package cluster

import "GoRedis/interface/resp"

// most commands that have a "key" structure similar to "set k1 v1"
// with the "key" in the second position can be executed using follow method.
func defaultFunc(d *Database, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	key := string(cmdArgs[1])         // the key in the second position
	peer := d.peerPicker.GetNode(key) // get peer node name by consistent hashing
	// send this command to target node,then receive the result and send it to the client . namely "replay"
	return d.relay(peer, c, cmdArgs)
}

func makeRouter() map[string]CmdFunc {
	routerMap := make(map[string]CmdFunc)
	routerMap["exists"] = defaultFunc
	routerMap["type"] = defaultFunc
	routerMap["set"] = defaultFunc
	routerMap["setnx"] = defaultFunc
	routerMap["get"] = defaultFunc
	routerMap["getset"] = defaultFunc
	routerMap["ping"] = ping
	routerMap["rename"] = rename
	routerMap["renamenx"] = rename
	routerMap["flushdb"] = flushdb
	routerMap["del"] = del
	routerMap["select"] = selectI
	return routerMap
}
