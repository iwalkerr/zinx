package ziface

// 链接管理模块
type IConnManager interface {
	Add(conn IConnection)                   // 添加连接
	Remove(conn IConnection)                // 删除连接
	Get(connId uint32) (IConnection, error) // 根据connId获取连接
	Len() int                               // 当前连接总数
	ClearConn()                             // 清除所有的连接
}
