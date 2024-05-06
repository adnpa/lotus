package ziface

type IConnManager interface {
	AddConn(conn IConnection)
	RemoveConn(conn IConnection)
	GetConn(connId uint32) (IConnection, error)
	GetConnNum() int
	ClearAllConn()
}
