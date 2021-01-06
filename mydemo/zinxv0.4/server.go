package main

import (
	"fmt"
	"gozinx/ziface"
	"gozinx/znet"
)

func main() {
	s := znet.NewServer("[zinx v0.4]")

	// 添加自定义router
	s.AddRouter((&PingRouter{}))

	s.Serve()
}

// 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

func (b *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("call handle PreHandle...")
	_, err := request.GetConnection().GetICPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back before ping error")
	}

}

// 主方法
func (b *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("call handle Handle...")
	_, err := request.GetConnection().GetICPConnection().Write([]byte("ping... ping...\n"))
	if err != nil {
		fmt.Println("call back ping... ping... error")
	}
}

func (b *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("call handle PostHandle...")
	_, err := request.GetConnection().GetICPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("call back after ping error")
	}
}
