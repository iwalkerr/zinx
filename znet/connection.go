package znet

import (
	"errors"
	"fmt"
	"gozinx/ziface"
	"io"
	"net"
)

type Connection struct {
	Conn       *net.TCPConn // 当前链接的socket TCP套接字
	ConnID     uint32       // 链接ID
	isClosed   bool         // 当前链接状态
	ExitChan   chan bool    // 告知当前链接已经退出/停止 channnel
	msgChan    chan []byte
	MsgHandler ziface.IMsgHandler // 该链接处理的方法
}

// 初始化链接模块方法
func NewConnection(conn *net.TCPConn, connID uint32, handler ziface.IMsgHandler) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: handler,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
	}
	return c
}

// 写消息
func (c *Connection) StartWriter() {
	fmt.Println("writer groutine is running")

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data error ", err)
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

// 链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("reader goroutine is running...")

	defer c.Stop()

	for {
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.GetICPConnection(), headData)
		if err != nil {
			fmt.Println("read msg head error")
			break
		}

		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}

		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetICPConnection(), data); err != nil {
				fmt.Println("read msg data error", err)
				break
			}
		}

		msg.SetData(data)

		req := Request{
			conn: c,
			msg:  msg,
		}

		go c.MsgHandler.DoMsgHandler(&req)
	}

}

// 发送数据， 将数据发送给远程的客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed when send msg")
	}

	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackege(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id= ", msgId)
		return errors.New("Pack error msg")
	}

	// 将数据发送给客户端
	c.msgChan <- binaryMsg

	return nil
}

// 启动链接
func (c *Connection) Start() {
	fmt.Println("conn start()... connId=", c.ConnID)

	go c.StartReader()
	go c.StartWriter()
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

	c.ExitChan <- true

	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
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
