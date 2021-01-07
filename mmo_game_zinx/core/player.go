package core

import (
	"fmt"
	"gozinx/ziface"
	"math/rand"
	"sync"

	"gozinx/mmo_game_zinx/pd"

	"google.golang.org/protobuf/proto"
)

type Player struct {
	Pid  int32
	Conn ziface.IConnection

	X float32
	Y float32
	Z float32
	V float32 // 旋转角度
}

var PidGen int32 = 1
var IdLock sync.Mutex

func NewPlayer(conn ziface.IConnection) *Player {
	IdLock.Lock()
	PidGen++
	IdLock.Unlock()

	p := &Player{
		Pid:  PidGen,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)),
		Y:    0,
		Z:    float32(140 + rand.Intn(20)),
		V:    0,
	}

	return p
}

func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("Marshal msg err:", err)
		return
	}

	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}

	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("Player SendMsg error!")
		return
	}
}

func (p *Player) SyncPid() {
	data := &pd.SyncPid{
		Pid: p.Pid,
	}
	p.SendMsg(1, data)
}

// 广播玩家
func (p *Player) BroadCastStartPostion() {
	msg := &pd.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pd.BroadCast_P{
			P: &pd.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	p.SendMsg(200, msg)
}

func (p *Player) Talk(content string) {
	msg := &pd.BroadCast{
		Pid: p.Pid,
		Tp:  1,
		Data: &pd.BroadCast_Content{
			Content: content,
		},
	}

	players := WorldMgrObj.GetAllPlayers()

	for _, player := range players {
		player.SendMsg(200, msg)
	}
}

// 同步聊天上线位置
func (p *Player) SyncSurrounding() {

	pids := WorldMgrObj.AoiMgr.GetPidByPos(p.X, p.Z)

	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}

	msg := &pd.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pd.BroadCast_P{
			P: &pd.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.Z,
			},
		},
	}

	for _, player := range players {
		player.SendMsg(200, msg)
	}

	playersMsg := make([]*pd.Player, 0, len(players))
	for _, player := range players {
		p := &pd.Player{
			Pid: player.Pid,
			P: &pd.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		playersMsg = append(playersMsg, p)
	}

	syncPlayerMsg := &pd.SyncPlayers{
		Ps: playersMsg[:],
	}

	p.SendMsg(202, syncPlayerMsg)
}

func (p *Player) UpdatePos(x, y, z, v float32) {
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	msg := &pd.BroadCast{
		Pid: p.Pid,
		Tp:  4,
		Data: &pd.BroadCast_P{
			P: &pd.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	players := p.GetSuroundingPlayers()

	for _, player := range players {
		player.SendMsg(200, msg)
	}
}

func (p *Player) GetSuroundingPlayers() []*Player {
	pids := WorldMgrObj.AoiMgr.GetPidByPos(p.X, p.Z)

	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}
	return players
}

func (p *Player) Offline() {
	players := p.GetSuroundingPlayers()

	msg := &pd.SyncPid{
		Pid: p.Pid,
	}

	for _, player := range players {
		player.SendMsg(201, msg)
	}

	WorldMgrObj.AoiMgr.RemoveFromGridByPos(int(p.Pid), p.X, p.Z)
	WorldMgrObj.RemovePlayerByPid(p.Pid)
}
