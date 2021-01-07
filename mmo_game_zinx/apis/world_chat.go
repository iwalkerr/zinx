package apis

import (
	"fmt"
	"gozinx/mmo_game_zinx/core"
	"gozinx/mmo_game_zinx/pd"
	"gozinx/ziface"
	"gozinx/znet"

	"google.golang.org/protobuf/proto"
)

type WorldChatApi struct {
	znet.BaseRouter
}

func (w *WorldChatApi) Handle(request ziface.IRequest) {
	msg := &pd.Talk{}
	if err := proto.Unmarshal(request.GetData(), msg); err != nil {
		fmt.Println("Talk  Unmarshal error", err)
		return
	}

	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		return
	}

	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	player.Talk(msg.Content)
}
