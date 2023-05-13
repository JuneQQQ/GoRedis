package cluster

import (
	"GoRedis/interface/resp"
	"GoRedis/lib/logger"
	"GoRedis/lib/utils"
	"GoRedis/resp/client"
	"GoRedis/resp/reply"
	"context"
	"errors"
	"strconv"
	"strings"
)

func (d *Database) borrowPeerClient(peer string) (*client.Client, error) {
	pool, ok := d.peerConnection[peer]
	if !ok {
		return nil, errors.New("connection not found")
	}
	object, err := pool.BorrowObject(context.Background())
	if err != nil {
		return nil, err
	}
	cli, ok := object.(*client.Client)
	if !ok {
		return nil, errors.New("-Err wrong type")
	}
	return cli, err
}

func (d *Database) returnPeerClient(peer string, c *client.Client) error {
	pool, ok := d.peerConnection[peer]
	if !ok {
		return errors.New("connection not found")
	}
	return pool.ReturnObject(context.Background(), pool)
}

func (d *Database) relay(peer string, c resp.Connection, args [][]byte) resp.Reply {
	if peer == d.self {
		// no need to relay
		logger.Debug("no need relay,exec command locally")
		return d.db.Exec(c, args)
	}

	peerClient, err := d.borrowPeerClient(peer)
	if err != nil {
		logger.Error("peer connection pool has been exhausted")
		return reply.MakeErrReply(err.Error())
	}
	defer func() {
		_ = d.returnPeerClient(peer, peerClient)
	}()

	// start for debug
	strArr := make([]string, len(args))
	for i, v := range args {
		strArr[i] = string(v)
	}
	str := strings.Join(strArr, " ")

	logger.Debug("switch to " + peer + " , command is : '" + str + "'")
	// end for debug

	// select standalone
	peerClient.Send(utils.ToCmdLine("select", strconv.Itoa(c.GetDBIndex())))
	// send command and receive result
	return peerClient.Send(args)
}

func (d *Database) broadcast(c resp.Connection, args [][]byte) map[string]resp.Reply {
	results := make(map[string]resp.Reply)
	for _, node := range d.peers {
		result := d.relay(node, c, args)
		results[node] = result
	}
	return results
}
