package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/zconf"
	"zinx/ziface"
	"zinx/zpack"
)

type Connection struct {
	Conn       *net.TCPConn           // 实际TCP链接
	ConnID     uint32                 // 链接ID
	isClosed   bool                   // 链接状态
	ExitChan   chan bool              // 告知当前链接已经退出
	MsgHandler ziface.IMessageHandler // 注册路由
	MsgChan    chan []byte            // 读写协程通讯管道
	IServer    ziface.IServer         //

	PropertyMap     map[string]interface{}
	propertyMapLock sync.RWMutex
}

// exitchan作用 reader进程判断对端状态是否关闭 通过管道告知写进程

// 初始化连接方法
func NewConnection(svr ziface.IServer, conn *net.TCPConn, connId uint32, msgHandler ziface.IMessageHandler) ziface.IConnection {
	c := &Connection{
		Conn:        conn,
		ConnID:      connId,
		isClosed:    false,
		ExitChan:    make(chan bool),
		MsgChan:     make(chan []byte),
		MsgHandler:  msgHandler,
		IServer:     svr,
		PropertyMap: make(map[string]interface{}),
	}

	c.IServer.GetConnManagemer().AddConn(c)
	return c
}

// 启动链接
func (conn *Connection) Start() {
	go conn.StartReader()
	go conn.StartWriter()

	conn.IServer.CallOnConnStart(conn)
}

// 读业务
func (conn *Connection) StartReader() {
	defer conn.Stop()
	for {
		// buf := make([]byte, zconf.GGlobalObj.MaxPackageSize)
		// _, err := conn.Conn.Read(buf)
		// if err != nil {
		// 	fmt.Println("server receive err: ", err)
		// 	continue
		// }

		dp := zpack.NewDatapack()
		// 读header
		headBin := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(conn.GetTcpCOnnection(), headBin)
		if err != nil {
			fmt.Println("read head err")
			break
		}

		msg, err := dp.Unpack(headBin)
		if err != nil {
			fmt.Println("unpack err: ", err)
			return
		}

		// 读数据
		if msg.GetMsgLen() > 0 {
			data := make([]byte, msg.GetMsgLen())
			_, err := io.ReadFull(conn.GetTcpCOnnection(), data)
			if err != nil {
				fmt.Println("read head err")
				return
			}
			msg.SetData(data)

			// 包装为Request，分发路由
			req := &Request{
				conn: conn,
				msg:  msg,
			}

			//执行路由
			if zconf.GGlobalObj.WorkerPoolSize > 0 {
				conn.MsgHandler.SendMsgToTaskQueue(req)
			} else {
				go conn.MsgHandler.DoMsgHandler(req)
			}
		}

	}
}

func (conn *Connection) StartWriter() {

	for {
		select {
		case data := <-conn.MsgChan:
			conn.GetTcpCOnnection().Write(data)
		case <-conn.ExitChan:
			return
		}
	}

}

// 关闭链接
func (conn *Connection) Stop() {

	if conn.isClosed {
		return
	}
	conn.isClosed = true

	conn.IServer.CallOnConnStop(conn)

	conn.Conn.Close()

	conn.ExitChan <- true

	conn.IServer.GetConnManagemer().RemoveConn(conn)
	close(conn.ExitChan)
	close(conn.MsgChan)
}

// 发送数据 写入管道 写进程实际发送
func (conn *Connection) Send(msgId uint32, data []byte) error {
	if conn.isClosed {
		return errors.New("err: connection closed")
	}
	dp := zpack.NewDatapack()
	msg := zpack.NewMessage(msgId, data)
	dataPack, err := dp.Pack(msg)
	if err != nil {
		return errors.New("pack err")
	}

	// if _, err := conn.GetTcpCOnnection().Write(dataPack); err != nil {
	// 	return errors.New("send datapack err")
	// }

	conn.MsgChan <- dataPack

	return nil
}

// Gettter
func (conn *Connection) GetTcpCOnnection() *net.TCPConn {
	return conn.Conn
}
func (conn *Connection) GetConnID() uint32 {
	return conn.ConnID
}
func (conn *Connection) RemoteAddr() net.Addr {
	return conn.Conn.RemoteAddr()
}

func (conn *Connection) SetProperty(key string, val interface{}) {
	conn.propertyMapLock.Lock()
	defer conn.propertyMapLock.Unlock()

	conn.PropertyMap[key] = val
}
func (conn *Connection) GetProperty(key string) (interface{}, error) {
	conn.propertyMapLock.RLock()
	defer conn.propertyMapLock.RUnlock()

	val, ok := conn.PropertyMap[key]
	if ok {
		return val, nil
	} else {
		return nil, errors.New("no property")
	}
}

func (conn *Connection) RemoveProperty(key string) {
	conn.propertyMapLock.Lock()
	defer conn.propertyMapLock.Unlock()

	delete(conn.PropertyMap, key)
}
