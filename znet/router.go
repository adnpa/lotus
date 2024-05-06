package znet

import "github.com/adnpa/lotus/ziface"

type BaseRouter struct{}

// 处理conn业务前的钩子方法 hook
func (br *BaseRouter) PreHandler(req ziface.IRequest) {}

// 处理conn业务的钩子方法 hook
func (br *BaseRouter) Handler(req ziface.IRequest) {}

// 处理conn业务后的钩子方法 hook
func (br *BaseRouter) PostHandler(req ziface.IRequest) {}
