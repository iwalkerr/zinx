package main

import (
	"fmt"
	"gozinx/ziface"
	"gozinx/znet"
)

func main() {
	s := znet.NewServer()

	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	// 添加自定义router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	s.Serve()
}

func DoConnectionBegin(conn ziface.IConnection) {
	// if err := conn.SendMsg(202, []byte("DoConnectionBegin")); err != nil {
	// 	fmt.Println(err)
	// }

	conn.SetProperty("name", "李发发发")
	conn.SetProperty("host", "12223222")
}

func DoConnectionLost(conn ziface.IConnection) {
	if _, err := conn.GetProperty("name"); err == nil {

	}
}

// 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// 主方法
func (b *PingRouter) Handle(request ziface.IRequest) {
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
	err := request.GetConnection().SendMsg(201, []byte("hello...hello...."))
	if err != nil {
		fmt.Println(err)
	}
}
