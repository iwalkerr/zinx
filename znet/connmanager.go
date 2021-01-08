package znet

import (
	"errors"
	"fmt"
	"gozinx/ziface"
	"sync"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection // 管理连接集合
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// 添加连接
func (c *ConnManager) Add(conn ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	c.connections[conn.GetConnID()] = conn
}

// 删除连接
func (c *ConnManager) Remove(conn ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	delete(c.connections, conn.GetConnID())
}

// 根据connId获取连接
func (c *ConnManager) Get(connId uint32) (ziface.IConnection, error) {
	c.connLock.RLock()
	defer c.connLock.RUnlock()

	if conn, ok := c.connections[connId]; ok {
		return conn, nil
	}
	return nil, errors.New("connection not found")
}

// 当前连接总数
func (c *ConnManager) Len() int {
	return len(c.connections)
}

// 清除所有的连接
func (c *ConnManager) ClearConn() {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	for connId, conn := range c.connections {
		conn.Stop()
		delete(c.connections, connId)
	}

	fmt.Println("clear all connection success")
}
