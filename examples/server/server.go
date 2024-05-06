package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
	// "fmt"
)

func main() {
	server := znet.NewServer("[zinxV0.7]")

	server.SetOnConnStart(OnConnStart)
	server.SetOnConnStop(OnConnStop)

	server.AddRouter(1, &PingRouter{})
	server.Serve()
}

type PingRouter struct {
	znet.BaseRouter
}

func OnConnStart(conn ziface.IConnection) {
	fmt.Println("on conn start")
	conn.SetProperty("aa", "abc")
}

func OnConnStop(conn ziface.IConnection) {
	fmt.Println("on conn stop")
	res, _ := conn.GetProperty("aa")
	fmt.Println("get pro", res)
}

// 处理conn业务前的钩子方法 hook
func (pr *PingRouter) PreHandler(req ziface.IRequest) {
	fmt.Println("prehandle\n")
	fmt.Println("get msg: ", string(req.GetData()))
	req.GetConnection().Send(1, []byte("echo"))
}

// 处理conn业务的钩子方法 hook
func (pr *PingRouter) Handler(req ziface.IRequest) {
	fmt.Println("handle\n")
	fmt.Println("get msg: ", string(req.GetData()))
	req.GetConnection().Send(1, []byte("echo"))
}

// 处理conn业务后的钩子方法 hook
func (pr *PingRouter) PostHandler(req ziface.IRequest) {
	fmt.Println("post handle\n")
	fmt.Println("get msg: ", string(req.GetData()))
	req.GetConnection().Send(1, []byte("echo"))
}
