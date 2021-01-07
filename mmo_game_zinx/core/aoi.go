package core

import "fmt"

// 定义AOI
const (
	AOI_MIN_X  = 85
	AOI_MAX_X  = 410
	AOI_CNTS_X = 10
	AOI_MIN_Y  = 75
	AOI_MAX_Y  = 400
	AOI_CNTS_Y = 20
)

type AOIManager struct {
	MinX  int
	MaxX  int
	CntsX int // x方向格子数量
	MinY  int
	MaxY  int
	CntsY int // y方向格子数量
	grids map[int]*Grid
}

func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grids: make(map[int]*Grid),
	}

	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			gid := y*cntsX + x

			aoiMgr.grids[gid] = NewGrid(
				gid,
				aoiMgr.MinX+x*aoiMgr.gridWidh(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWidh(),
				aoiMgr.MinY+y*aoiMgr.gridLength(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridLength(),
			)
		}
	}

	return aoiMgr
}

// x轴宽度
func (m *AOIManager) gridWidh() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

// y轴宽度
func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

func (m *AOIManager) String() string {
	return fmt.Sprintf("%+v\n", *m)
}

func (m *AOIManager) GetSurroundGridsByGid(gId int) (gids []*Grid) {
	if _, ok := m.grids[gId]; !ok {
		return
	}

	gids = append(gids, m.grids[gId])

	idx := gId % m.CntsX

	if idx > 0 {
		gids = append(gids, m.grids[gId-1])
	}

	if idx < m.CntsX-1 {
		gids = append(gids, m.grids[gId+1])
	}

	gidx := make([]int, 0, len(gids))
	for _, v := range gids {
		gidx = append(gidx, v.GID)
	}

	for _, v := range gidx {
		idy := v / m.CntsY

		if idy > 0 {
			gids = append(gids, m.grids[v-m.CntsX])
		}

		if idy < m.CntsY-1 {
			gids = append(gids, m.grids[v+m.CntsX])
		}
	}
	return gids
}

func (m *AOIManager) GetGidByPos(x, y float32) int {
	idx := (int(x) - m.MinX) / m.gridWidh()
	idy := (int(y) - m.MinY) / m.gridLength()

	return idy*m.CntsX + idx
}

func (m *AOIManager) GetPidByPos(x, y float32) (playerIds []int) {
	gid := m.GetGidByPos(x, y)

	gids := m.GetSurroundGridsByGid(gid)

	for _, grid := range gids {
		playerIds = append(playerIds, grid.GetPlayerIds()...)
	}
	return
}

func (m *AOIManager) AddPidToGrid(pId, gId int) {
	m.grids[gId].Add(pId)
}

func (m *AOIManager) RemovePidFromGrid(pId, gId int) {
	m.grids[gId].Remove(pId)
}

func (m *AOIManager) GetPidsByGid(gId int) (playerIds []int) {
	playerIds = m.grids[gId].GetPlayerIds()
	return
}

func (m *AOIManager) AddToGridByPos(pId int, x, y float32) {
	gId := m.GetGidByPos(x, y)
	grid := m.grids[gId]
	grid.Add(pId)
}

func (m *AOIManager) RemoveFromGridByPos(pId int, x, y float32) {
	gId := m.GetGidByPos(x, y)
	grid := m.grids[gId]

	grid.Remove(gId)
}
