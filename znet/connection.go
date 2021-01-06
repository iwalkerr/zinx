package znet

import (
	"fmt"
	"gozinx/utils"
	"gozinx/ziface"
	"net"
)

type Connection struct {
	Conn     *net.TCPConn   // 当前链接的socket TCP套接字
	ConnID   uint32         // 链接ID
	isClosed bool           // 当前链接状态
	ExitChan chan bool      // 告知当前链接已经退出/停止 channnel
	Router   ziface.IRouter // 该链接处理的方法
}

// 初始化链接模块方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}
	return c
}

// 链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("reader goroutine is running...")

	defer c.Stop()

	for {
		// 读取客户端到buf中
		buf := make([]byte, int(utils.GlobalObject.MaxPackageSize))
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err", err)
			continue
		}

		req := Request{
			conn: c,
			data: buf,
		}

		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

	}

}

// 启动链接
func (c *Connection) Start() {
	fmt.Println("conn start()... connId=", c.ConnID)

	go c.StartReader()
}

// 停止链接
func (c *Connection) Stop() {
	fmt.Println("conn stop()... connID=", c.ConnID)

	if c.isClosed {
		return
	}
	c.isClosed = true

	// 关闭socket链接
	c.Conn.Close()
	close(c.ExitChan)
}

// 获取当前链接绑定的socket conn
func (c *Connection) GetICPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接模块ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据， 将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil
}
