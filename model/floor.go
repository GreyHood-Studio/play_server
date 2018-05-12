package model

import (
	"github.com/gin-gonic/gin/json"
	"github.com/GreyHood-Studio/server_util/error"
)

// 한 층의 데이터 일반적인 게임에서의 room 개념과 유사
type Floor struct {
	Status		FloorStatus	// 맵에 대한 현재 상태 구조체
	Players		map[int]Player	// 맵에 존재하는 플레이어들의 리스트 ( user, enemy, character )
}

// 한 층의 현재 상태 ( 동접, 방 관리를 위한 내용 )
// 플레이어 수, ai수, etc...
type FloorStatus struct {
	PlayerCount		int
	MapType			int
}

func (floor *Floor) GetFloorStatus() []byte {
	jsonBytes, err := json.Marshal(floor)
	error.CheckError(err,"error get floor status")
	return jsonBytes
}

func (floor *Floor) HandleFloor() {

}

func (floor *Floor) AddPlayer(playerName string, playerId int) {
	player := Player{PlayerName:playerName}

	floor.Players[playerId] = player
}

func (floor *Floor) DeletePlayer(playerId int) {
	delete(floor.Players, playerId)
}

func NewFloor() *Floor {
	floor := &Floor{
		Players: make(map[int]Player),
	}
	return floor
}