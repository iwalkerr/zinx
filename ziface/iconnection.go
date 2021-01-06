package ziface

import "net"

type IConnection interface {
	Start()                         // 启动链接
	Stop()                          // 停止链接
	GetICPConnection() *net.TCPConn // 获取当前链接绑定的socket conn
	GetConnID() uint32              // 获取当前连接模块ID
	RemoteAddr() net.Addr           // 获取远程客户端的TCP状态 IP port
	Send(data []byte) error         // 发送数据， 将数据发送给远程的客户端
}

// 定义一个处理链接业务的方法
type HanleFunc func(*net.TCPConn, []byte, int) error
