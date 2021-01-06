package ziface

type IRequest interface {
	GetConnection() IConnection // 得到当前链接
	GetData() []byte            // 得到请求的消息数据
	GetMsgID() uint32           // 消息ID
}
