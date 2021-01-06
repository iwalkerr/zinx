package main

import (
	"fmt"
	"gozinx/ziface"
	"gozinx/znet"
)

func main() {
	s := znet.NewServer("[zinx v0.5]")

	// 添加自定义router
	s.AddRouter((&PingRouter{}))

	s.Serve()
}

// 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// 主方法
func (b *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("call handle Handle...")

	err := request.GetConnection().SendMsg(1, []byte("ping...ping...."))
	if err != nil {
		fmt.Println(err)
	}
}
