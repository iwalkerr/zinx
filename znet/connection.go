package znet

import (
	"fmt"
	"gozinx/ziface"
	"net"
)

type Connection struct {
	Conn      *net.TCPConn     // 当前链接的socket TCP套接字
	ConnID    uint32           // 链接ID
	isClosed  bool             // 当前链接状态
	handleAPI ziface.HanleFunc // 当前链接所绑定的处理业务方法API
	ExitChan  chan bool        // 告知当前链接已经退出/停止 channnel
}

// 初始化链接模块方法
func NewConnection(conn *net.TCPConn, connID uint32, callbackApi ziface.HanleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		handleAPI: callbackApi,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
	}
	return c
}

// 链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("reader goroutine is running...")

	defer c.Stop()

	for {
		// 读取客户端到buf中，最大512字节
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err", err)
			continue
		}

		// 调用链接所绑定的handAPI
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("connID ", c.ConnID, " handle is error", err)
			break
		}
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
func (c *Connection) Send() error {
	return nil
}
