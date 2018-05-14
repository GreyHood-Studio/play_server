package model

import (
	"github.com/gin-gonic/gin/json"
	"github.com/GreyHood-Studio/server_util/error"
)

// 한 층의 데이터 일반적인 게임에서의 room 개념과 유사
type Room struct {
	Status		RoomStatus	// 맵에 대한 현재 상태 구조체
	Players		map[int]Player	// 맵에 존재하는 플레이어들의 리스트 ( user, enemy, character )
}

// 한 층의 현재 상태 ( 동접, 방 관리를 위한 내용 )
// 플레이어 수, ai수, etc...
type RoomStatus struct {
	PlayerCount		int
	MapType			int
}

func (room *Room) GetRoomStatus() []byte {
	jsonBytes, err := json.Marshal(room)
	error.CheckError(err,"error get room status")
	return jsonBytes
}

func (room *Room) HandleRoom() {

}

func (room *Room) AddPlayer(playerName string, playerId int) {
	player := Player{PlayerName:playerName, PlayerId:playerId}

	room.Players[playerId] = player
}

func (room *Room) DeletePlayer(playerId int) {
	delete(room.Players, playerId)
}

func NewRoom() *Room {
	room := &Room{
		Players: make(map[int]Player),
	}
	return room
}