package tcp

import (
	"GoRedis/config"
	"GoRedis/interface/tcp"
	"GoRedis/lib/logger"
	"context"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

type Config struct {
	Address string
}

func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	listener, err := net.Listen("tcp", cfg.Address)
	closeChan := make(chan struct{})

	// 注册 signalChan 关心的些系统信号，此chan将关闭信号间接传递给了 closeChan
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		switch <-signalChan {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()

	if err != nil {
		return err
	}
	logger.Info("server started in " + config.Properties.Bind + ":" + strconv.Itoa(config.Properties.Port))

	ListenAndServe(listener, handler, closeChan)

	return nil
}

func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {
	go func() {
		<-closeChan
		logger.Info("shutting down...")
		_ = listener.Close()
		_ = handler.Close()
	}()

	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()

	ctx := context.Background()
	var waitDone sync.WaitGroup
	for true {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("accept error,msg:", err.Error())
			break
		}
		logger.Info("accepted link,RemoteAddr ", conn.RemoteAddr())

		go func() {
			defer func() {
				waitDone.Done()
			}()
			waitDone.Add(1)
			handler.Handle(ctx, conn)
		}()
	}
	waitDone.Wait()
}
