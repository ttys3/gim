package config

import (
	"os"
)

var (
	Logic    LogicConf
	Conn     ConnConf
	Business BusinessConf
)

// logic配置
type LogicConf struct {
	MySQL            string
	NSQIP            string
	RedisIP          string
	RedisPassword    string
	RPCIntListenAddr string
	RPCExtListenAddr string
	ConnRPCAddrs     string
	BusinessRPCAddrs string
}

// TCPConnConf配置
type ConnectConf struct {
	TCPListenAddr int
	WSListenAddr  string
	RPCListenAddr string
	LocalAddr     string
	LogicRPCAddrs string
}

// Business配置
type BusinessConf struct {
	MySQL            string
	NSQIP            string
	RedisIP          string
	RedisPassword    string
	RPCIntListenAddr string
	RPCExtListenAddr string
	LogicRPCAddrs    string
}

func init() {
	env := os.Getenv("gim_env")
	switch env {
	case "dev":
		initDevConf()
	case "prod":
		initProdConf()
	default:
		initLocalConf()
	}
}
