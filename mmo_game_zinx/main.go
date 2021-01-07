package main

import (
	"gozinx/mmo_game_zinx/apis"
	"gozinx/mmo_game_zinx/core"
	"gozinx/ziface"
	"gozinx/znet"
)

func main() {
	s := znet.NewServer("MMO Game Zinx")

	// hook
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)

	// 注册路由
	s.AddRouter(2, &apis.WorldChatApi{})
	// 移动路由
	s.AddRouter(3, &apis.MoveApi{})

	s.Serve()
}

func OnConnectionLost(conn ziface.IConnection) {
	pid, _ := conn.GetProperty("pid")

	payer := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	// 下线
	payer.Offline()
}

func OnConnectionAdd(conn ziface.IConnection) {
	player := core.NewPlayer(conn)

	// 同步pid
	player.SyncPid()

	// 同步初始位置
	player.BroadCastStartPostion()

	// 玩家添加到world
	core.WorldMgrObj.AddPlayer(player)

	// 添加pid
	conn.SetProperty("pid", player.Pid)

	// 同步周边玩家，广播当前位置信息
	player.SyncSurrounding()
}
