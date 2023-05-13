package tcp

import (
	"GoRedis/lib/logger"
	"GoRedis/lib/sync/wait"
	"bufio"
	"context"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// EchoClient Client information for the task assigned to this handler
type EchoClient struct {
	Conn     net.Conn
	Handling wait.Wait // Number of goroutines processing core tasks,cannot close the connection
}

// Close the connection to this client
func (client *EchoClient) Close() error {
	// Wait for three seconds at the latest to close the connection
	client.Handling.WaitWithTimeout(3 * time.Second)
	_ = client.Conn.Close()
	return nil
}

type EchoHandler struct {
	activeConn sync.Map    // all relevant client information
	closed     atomic.Bool // whether the handler has been closed
}

func (handler *EchoHandler) Handle(ctx context.Context, conn net.Conn) {

	if handler.closed.Load() {
		_ = conn.Close()
	}

	client := &EchoClient{
		Conn: conn,
	}

	// +1 means the task for this client by the server is in progress,
	// cannot close client,unless it times out
	client.Handling.Add(1)

	defer func() {
		client.Handling.Done()
	}()

	handler.activeConn.Store(client, struct{}{})

	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("Connection Closed,RemoteAddr", conn.RemoteAddr())
				handler.activeConn.Delete(client)
			} else {
				logger.Error("read error,msg:", err.Error())
			}
			return
		}

		b := []byte(msg)
		_, _ = conn.Write(b)
	}
}

// Close this handler, all ongoing tasks assigned to this handler will be terminated,
// because the connection was closed
func (handler *EchoHandler) Close() error {
	logger.Info("handler shutting down")
	handler.closed.Store(true)
	handler.activeConn.Range(func(key, value interface{}) bool {
		client := key.(*EchoClient)
		_ = client.Conn.Close()
		return true
	})
	return nil
}

func MakeEchoHandler() *EchoHandler {
	return &EchoHandler{}
}
