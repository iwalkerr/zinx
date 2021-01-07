package apis

import (
	"fmt"
	"gozinx/mmo_game_zinx/core"
	"gozinx/mmo_game_zinx/pd"
	"gozinx/ziface"
	"gozinx/znet"

	"google.golang.org/protobuf/proto"
)

type MoveApi struct {
	znet.BaseRouter
}

func (w *MoveApi) Handle(request ziface.IRequest) {
	msg := &pd.Position{}

	if err := proto.Unmarshal(request.GetData(), msg); err != nil {
		fmt.Println("Move: Position Unmarshal error", err)
		return
	}

	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid error ", err)
		return
	}

	fmt.Println("pid ", pid)

	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	player.UpdatePos(msg.X, msg.Y, msg.Z, msg.V)
}
