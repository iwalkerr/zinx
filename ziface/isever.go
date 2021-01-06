package ziface

// 定义一个服务器端口
type IServer interface {
	Start()
	Stop()
	Serve()
}
