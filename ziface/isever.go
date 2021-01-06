package ziface

// 定义一个服务器端口
type IServer interface {
	Start()                                 // 启动服务
	Stop()                                  // 停止服务
	Serve()                                 // 运行服务
	AddRouter(msgId uint32, router IRouter) // 注册路由功能
}
