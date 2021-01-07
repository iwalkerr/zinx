package core

import "sync"

// 世界的管理模块

type WorldManager struct {
	AoiMgr  *AOIManager
	Players map[int32]*Player
	pLock   sync.RWMutex
}

// 提供一个对外世界的管理模块
var WorldMgrObj *WorldManager

func init() {
	WorldMgrObj = &WorldManager{
		AoiMgr:  NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y),
		Players: make(map[int32]*Player),
	}
}

func (m *WorldManager) AddPlayer(player *Player) {
	m.pLock.Lock()
	m.Players[player.Pid] = player
	m.pLock.Unlock()

	m.AoiMgr.AddToGridByPos(int(player.Pid), player.X, player.Z)
}

func (m *WorldManager) RemovePlayerByPid(pid int32) {
	player := m.Players[pid]
	m.AoiMgr.RemoveFromGridByPos(int(pid), player.X, player.Z)

	m.pLock.Lock()
	delete(m.Players, pid)
	m.pLock.Unlock()
}

func (m *WorldManager) GetPlayerByPid(pid int32) *Player {
	m.pLock.RLock()
	defer m.pLock.RUnlock()

	return m.Players[pid]
}

func (m *WorldManager) GetAllPlayers() []*Player {
	m.pLock.RLock()
	defer m.pLock.RUnlock()

	players := make([]*Player, 0)

	for _, p := range m.Players {
		players = append(players, p)
	}

	return players
}
