package ziface

// 客户端抽象为connection

import "net"

type IConnection interface {
	Start()
	Stop()
	GetTcpCOnnection() *net.TCPConn
	GetConnID() uint32
	RemoteAddr() net.Addr
	Send(uint32, []byte) error

	SetProperty(key string, val interface{})
	GetProperty(key string)(interface{}, error)
	RemoveProperty(key string)
}

// type HandleFunc func(*net.TCPConn, []byte, int) error
