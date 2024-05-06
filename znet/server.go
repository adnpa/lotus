package znet

import (
	"fmt"
	"net"
	"zinx/zconf"
	"zinx/ziface"
)

type Server struct {
	Name        string
	IPVersion   string
	Ip          string
	Port        int
	MsgHandler  ziface.IMessageHandler
	ConnManager ziface.IConnManager

	OnConnStart func(conn ziface.IConnection)
	OnConnStop  func(conn ziface.IConnection)
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:        zconf.GGlobalObj.Name,
		IPVersion:   "tcp4",
		Ip:          zconf.GGlobalObj.Host,
		Port:        zconf.GGlobalObj.TcpPort,
		MsgHandler:  NewMessageHandler(),
		ConnManager: NewConnManager(),
	}
	return s
}

func (s *Server) Start() {
	go s.MsgHandler.StartWorkerPool()

	// 1. resolve
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("init socket errl:", err)
		return
	}
	// 2. listen
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("listen err:", err)
		return
	}
	defer listener.Close()

	// 3. accept & handler msg
	var cid uint32 = 0
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("accept err: ", err)
			continue
		}

		// 如果大于最大连接 拒绝服务
		if s.ConnManager.GetConnNum() >= zconf.GGlobalObj.MaxConn {
			conn.Close()
			continue
		}

		clientConn := NewConnection(s, conn, cid, s.MsgHandler)

		clientConn.Start()
		cid++
	}
}

func (s *Server) Stop() {
	s.ConnManager.ClearAllConn()
}

func (s *Server) Serve() {
	go s.Start()
	fmt.Println("server start ok...")

	// todo addition work

	select {}
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("add router succ")
}

func (s *Server) GetConnManagemer() ziface.IConnManager {
	return s.ConnManager
}

func (s *Server) SetOnConnStart(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		s.OnConnStop(conn)
	}
}

// 业务处理 回调方法
// func CallbackClient(conn *net.TCPConn, buf []byte, cnt int) error {
// 	fmt.Println("callback handle")
// 	if _, err := conn.Write(buf[:cnt]); err != nil {
// 		fmt.Println("write err", err)
// 		return errors.New("CallBackToClient error")
// 	}
// 	return nil
// }
