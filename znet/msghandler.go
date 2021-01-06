package znet

import (
	"fmt"
	"gozinx/ziface"
)

type MsgHandler struct {
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

func (m *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	handler, ok := m.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId =", request.GetMsgID(), "is not found")
		return
	}

	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (m *MsgHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	if _, ok := m.Apis[msgId]; ok {
		//id已经注册
		fmt.Println("repeat api,msgID=", msgId)
		return
	}

	m.Apis[msgId] = router
	fmt.Println("add api msgId ", msgId)
}
