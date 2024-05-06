package znet

import (
	"fmt"
	"zinx/zconf"
	"zinx/ziface"
)

type MessageHandler struct {
	IdRouterMap    map[uint32]ziface.IRouter
	TaskQueue      []chan ziface.IRequest
	WorkerPoolSize uint32
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		IdRouterMap:    make(map[uint32]ziface.IRouter),
		TaskQueue:      make([]chan ziface.IRequest, zconf.GGlobalObj.WorkerPoolSize),
		WorkerPoolSize: zconf.GGlobalObj.WorkerPoolSize,
	}
}

func (mh *MessageHandler) DoMsgHandler(req ziface.IRequest) {
	handler, ok := mh.IdRouterMap[req.GetMsgId()]
	if !ok {
		fmt.Println("no router")
	}
	handler.PreHandler(req)
	handler.Handler(req)
	handler.PostHandler(req)
}

func (mh *MessageHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	if _, ok := mh.IdRouterMap[msgId]; ok {
		panic("repeat api")
	}
	mh.IdRouterMap[msgId] = router
}

func (mh *MessageHandler) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan ziface.IRequest, zconf.GGlobalObj.MaxWorkerTaskLen)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}
func (mh *MessageHandler) StartOneWorker(workerId int, taskQueue chan ziface.IRequest) {

	for req := range taskQueue {
		// select {
		// case req := <-taskQueue:
		fmt.Printf("worker %d handle message\n", workerId)
		mh.DoMsgHandler(req)
		// }
	}
}

// 发到到worker对应channel
func (mh *MessageHandler) SendMsgToTaskQueue(req ziface.IRequest) {
	fmt.Println("get req")

	// 负载均衡 选择worker
	workerId := req.GetConnection().GetConnID() % mh.WorkerPoolSize

	// 发到到worker对应channel
	mh.TaskQueue[workerId] <- req
}
