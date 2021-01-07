package ziface

type IMsgHandler interface {
	DoMsgHandler(request IRequest)
	AddRouter(msgId uint32, router IRouter)
	StartWorkerPool()                    // 启动worker工作池
	SendMsgToTaskQueue(request IRequest) // 将消息发送给消息任务队列处理
}
