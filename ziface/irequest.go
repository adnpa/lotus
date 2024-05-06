package ziface

// 链接和请求数据封装为Request

type IRequest interface {
	GetConnection() IConnection
	GetData() []byte
	GetMsgId() uint32
}
 