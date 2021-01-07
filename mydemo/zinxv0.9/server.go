package main

import (
	"fmt"
	"gozinx/ziface"
	"gozinx/znet"
)

func main() {
	s := znet.NewServer("[zinx v0.8]")

	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	// 添加自定义router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	s.Serve()
}

func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("==> DoConnectionBegin")
	if err := conn.SendMsg(202, []byte("DoConnectionBegin")); err != nil {
		fmt.Println(err)
	}
}

func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("===> 回收资源")
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

// 自定义路由
type HelloRouter struct {
	znet.BaseRouter
}

// 主方法
func (b *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("call HelloRouter...")

	err := request.GetConnection().SendMsg(201, []byte("hello...hello...."))
	if err != nil {
		fmt.Println(err)
	}
}
