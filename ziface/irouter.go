package ziface

// 路由抽象接口
type IRouter interface {
	// 在处理conn业务之前的钩子方法
	PreHandle(request IRequest)
	// 在处理conn业务主方法
	Handle(request IRequest)
	// 在处理conn之后的钩子方法
	PostHandle(request IRequest)
}
