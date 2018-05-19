package model

import (
	"github.com/gin-gonic/gin/json"
	"github.com/GreyHood-Studio/server_util/checker"
)

// 한 층의 데이터 일반적인 게임에서의 room 개념과 유사
type Room struct {
	// 맵에 대한 정보
	MapId			int
	MaxUser			int
	// 맵에 존재하는 플레이어들의 리스트 ( user, enemy, character )
	Players			map[int]*Player
	// Map에 존재하는 Player들의 리스트
	MapObjects		[]MapObject
	// monster, npc, histories
	LastPlayerNum	int
	LastBulletNum	int

	Bullets			map[int]Bullet
}

func (room *Room) GetRoomStatus() []byte {
	jsonBytes, err := json.Marshal(room)
	checker.CheckError(err,"error get room status")
	return jsonBytes
}

func (room *Room) AddPlayer(newPlayer *Player) int {
	newPlayerNum := room.LastPlayerNum
	newPlayer.PlayerNum = newPlayerNum
	room.Players[newPlayerNum] = newPlayer

	room.LastPlayerNum++
	//golang이 last player num ++가 안되서, 추가하고 삭제
	return room.LastPlayerNum - 1
}

func (room *Room) DeletePlayer(playerNum int) {
	delete(room.Players, playerNum)
}

func NewRoom(maxUser int, mapId int) *Room {
	room := &Room{
		MapId: mapId,
		MaxUser: maxUser,
		Players: make(map[int]*Player),
		MapObjects: make([]MapObject, 100),
		Bullets: make(map[int]Bullet),
	}
	return room
}