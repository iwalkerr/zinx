package znet

import (
	"fmt"
	"gozinx/utils"
	"gozinx/ziface"
)

type MsgHandler struct {
	Apis           map[uint32]ziface.IRouter
	TaskQueue      []chan ziface.IRequest // 消息队列 - 取任务
	WorkerPoolSize uint32                 // worker 池
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// 处理消息
func (m *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	handler, ok := m.Apis[request.GetMsgID()]
	if !ok {
		fmt.Printf("api msgId = %d is not found\n", int(request.GetMsgID()))
		return
	}

	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 添加路由
func (m *MsgHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	if _, ok := m.Apis[msgId]; ok {
		//id已经注册
		fmt.Println("repeat api,msgID=", msgId)
		return
	}

	m.Apis[msgId] = router
}

// 启动一个worker工作池
func (m *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkTaskLen)
		go m.startOneWorker(i, m.TaskQueue[i])
	}
}

// 启动一个工作流程
func (m *MsgHandler) startOneWorker(workerId int, taskQueue chan ziface.IRequest) {
	for {
		request := <-taskQueue
		m.DoMsgHandler(request)
	}
}

func (m *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	workerId := request.GetConnection().GetConnID() % m.WorkerPoolSize
	// fmt.Println("request msgId to workerId =", workerId)

	m.TaskQueue[workerId] <- request
}
