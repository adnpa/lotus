package ziface

// router模块，负责分发不同处理方法，例如不同消息对应不同处理方式
// 模板方法

type IRouter interface {
	PreHandler(req IRequest)
	Handler(req IRequest)
	PostHandler(req IRequest)
}
