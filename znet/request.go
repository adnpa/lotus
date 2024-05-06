package znet

import "zinx/ziface"

type Request struct {
	conn ziface.IConnection
	msg  ziface.IMessage
}

func (req *Request) GetConnection() ziface.IConnection {
	return req.conn
}

func (req *Request) GetData() []byte {
	return req.msg.GetData()
}

func (req *Request) GetMsgId() uint32 {
	return req.msg.GetMsgId()
}
