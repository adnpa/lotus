package ziface

type IServer interface {
	Start()
	Stop()
	Serve()
	AddRouter(uint32, IRouter)
	GetConnManagemer() IConnManager

	SetOnConnStart(func(conn IConnection))
	SetOnConnStop(func(conn IConnection))
	CallOnConnStart(conn IConnection)
	CallOnConnStop(conn IConnection)
}
