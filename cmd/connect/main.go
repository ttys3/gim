package main

import (
	"context"
	"gim/config"
	"gim/internal/connect"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/rpc"
	"gim/pkg/util"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	logger.Init()
	// 启动rpc服务
	go func() {
		defer util.RecoverPanic()
		connect.StartRPCServer()
	}()

	// 初始化Rpc Client
	rpc.InitLogicIntClient(config.Conn.LogicRPCAddrs)

	// 启动TCP长链接服务器
	go func() {
		connect.StartTCPServer()
	}()

	// 启动WebSocket长链接服务器
	go func() {
		defer util.RecoverPanic()
		connect.StartWSServer(config.Conn.WSListenAddr)
	}()

	c := make(chan os.Signal, 0)
	signal.Notify(c, syscall.SIGTERM)

	s := <-c
	logger.Logger.Info("server stop start", zap.Any("signal", s))
	rpc.LogicIntClient.ServerStop(context.TODO(), &pb.ServerStopReq{ConnAddr: config.Conn.LocalAddr})
	logger.Logger.Info("server stop end")
}
