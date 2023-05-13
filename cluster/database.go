package cluster

import (
	"GoRedis/config"
	"GoRedis/interface/database"
	"GoRedis/interface/resp"
	"GoRedis/lib/consistenthashing"
	"GoRedis/lib/logger"
	"GoRedis/resp/reply"
	database2 "GoRedis/standalone"
	"context"
	pool "github.com/jolestar/go-commons-pool/v2"
	"strings"
)

type Database struct {
	self           string                      // self node name
	peers          []string                    // peer node names
	peerPicker     *consistenthashing.NodeMap  // a struct for consistent hashing
	peerConnection map[string]*pool.ObjectPool // connection pool of peer nodes
	db             database.IDatabase          // local storage
}

var routerMap = makeRouter() // command name -> method

func MakeClusterDatabase() *Database {
	cluster := &Database{
		self:           config.Properties.Self,
		db:             database2.NewStandaloneDatabase(),
		peerPicker:     consistenthashing.NewNodeMap(nil),
		peerConnection: make(map[string]*pool.ObjectPool),
		peers:          make([]string, 0, len(config.Properties.Peers)+1),
	}

	// init trivial properties
	for _, peer := range config.Properties.Peers {
		cluster.peers = append(cluster.peers, peer)
	}
	cluster.peers = append(cluster.peers, config.Properties.Self)

	cluster.peerPicker.AddNode(cluster.peers...)

	for _, peer := range config.Properties.Peers {
		cluster.peerConnection[peer] =
			pool.NewObjectPoolWithDefaultConfig(context.Background(), &connectionFactory{Peer: peer})
		cluster.peerConnection[peer].Config.MaxTotal = 100
	}
	return cluster
}

func (d *Database) Exec(conn resp.Connection, args database.CmdLine) (result resp.Reply) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
			result = reply.MakeUnknownErrReply()
		}
	}()

	cmdFunc, ok := routerMap[strings.ToLower(string(args[0]))]

	if !ok {
		return reply.MakeErrReply("Err not supported")
	}

	result = cmdFunc(d, conn, args)

	return
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) AfterClientClose(c resp.Connection) {
	d.db.AfterClientClose(c)
}

type CmdFunc func(d *Database, c resp.Connection, cmdArgs [][]byte) resp.Reply
