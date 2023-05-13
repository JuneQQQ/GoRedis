package cluster

import (
	"GoRedis/resp/client"
	"context"
	"errors"
	pool "github.com/jolestar/go-commons-pool/v2"
)

type connectionFactory struct {
	Peer string
}

func (conn *connectionFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	cli, err := client.MakeClient(conn.Peer)
	if err != nil {
		return nil, err
	}
	cli.Start()
	return pool.NewPooledObject(cli), nil
}

func (conn *connectionFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	c, ok := object.Object.(*client.Client)
	if !ok {
		return errors.New("type dismatch")
	}
	c.Close()
	return nil
}

func (conn *connectionFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	return true
}

func (conn *connectionFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}

func (conn *connectionFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}
