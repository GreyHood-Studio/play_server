package model

type FloorMap struct {
	// Room들의 리스트
	rooms		[]Room
	// Door들의 리스트
	doors		[]Pos
	x_size		int
	y_size		int
}

type Room struct {
	tiles		[][]tile
	x_start		int
	y_start		int
}

// 타일에 대한 정보
type tile struct {
	tileType		int
}

func (r *Room) changeTypeType(x int, y int, tileType int) {
	r.tiles[x][y].tileType = tileType
}

func (r *Room) getRoomSize() int {
	return cap(r.tiles) * cap(r.tiles[0])
}

func (m *FloorMap) getTotalMapSize() int {
	// 이게 더 빠를꺼 같긴함
	return m.x_size * m.y_size
	//
}
