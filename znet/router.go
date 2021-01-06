package znet

import "gozinx/ziface"

type BaseRouter struct{}

// 之前
func (b *BaseRouter) PreHandle(request ziface.IRequest) {

}

// 主方法
func (b *BaseRouter) Handle(request ziface.IRequest) {

}

// 之后
func (b *BaseRouter) PostHandle(request ziface.IRequest) {

}
