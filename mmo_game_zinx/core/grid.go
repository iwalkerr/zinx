package core

import (
	"fmt"
	"sync"
)

type Grid struct {
	GID       int
	MinX      int
	MaxX      int
	MinY      int
	MaxY      int
	playerIds map[int]bool
	pIdLock   sync.RWMutex
}

func NewGrid(gId, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gId,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIds: make(map[int]bool),
	}
}

func (g *Grid) Add(playerId int) {
	g.pIdLock.Lock()
	defer g.pIdLock.Unlock()

	g.playerIds[playerId] = true
}

func (g *Grid) Remove(playerId int) {
	g.pIdLock.Lock()
	defer g.pIdLock.Unlock()

	delete(g.playerIds, playerId)
}

func (g *Grid) GetPlayerIds() (playerIds []int) {
	g.pIdLock.RLock()
	defer g.pIdLock.RUnlock()

	for k := range g.playerIds {
		playerIds = append(playerIds, k)
	}
	return
}

func (g *Grid) String() string {
	return fmt.Sprintf("%+v\n", &g)
}
