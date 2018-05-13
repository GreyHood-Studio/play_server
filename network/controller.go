package network

import (
	"github.com/GreyHood-Studio/play_server/network/protocol"
	"fmt"
)

// key = server_id, value: bullet_id
var bulletID map[int]int
var playerID map[int]int

func init()  {
	bulletID = make(map[int]int)
	playerID = make(map[int]int)
}

func (gameClient *gameClient) handleCommon(event protocol.CommonEvent) {
	switch event.EventType {
	case 1: // firebullet
	case 2: // hitbullet
	case 3: // reloadweapon
	case 4: // evade
	case 5: // deadplayer
	case 6: // exitplayer
	}
}

func (gameClient *gameClient) handleBullet(bulletBytes []byte) []byte{
	bulletID[gameClient.serverId]++
	bulletPacket := protocol.AssignBulletID(bulletID[gameClient.serverId],bulletBytes)
	return bulletPacket
}

func (gameClient *gameClient) handleStart(startBytes []byte) ([]byte, []byte){
	startEvent := protocol.UnpackStart(startBytes)
	players := floorMap[gameClient.serverId].Players
	for _, value := range players {
		// 기존 유저의 리스트를 받아옴
		fmt.Printf("current user[%d]: %s",value.PlayerId, value.PlayerName)
		playerEvent := protocol.PlayerEvent{ PlayerName: value.PlayerName, PlayerId: value.PlayerId }
		startEvent.PlayerList = append(startEvent.PlayerList, playerEvent)
	}

	// game client에 client id 할당
	gameClient.clientId = playerID[gameClient.serverId]
	// game floor의 player list에 id 할당
	floorMap[gameClient.serverId].AddPlayer(gameClient.clientName, gameClient.clientId)
	fmt.Println("assign new player id[",playerID[gameClient.serverId],"]")
	// 보낼 데이터에 id랑 player name 할당
	startEvent.PlayerList[0].PlayerId = playerID[gameClient.serverId]
	startEvent.PlayerList[0].PlayerName = gameClient.clientName
	// 다 한 뒤에 다음 client server id 증가
	playerID[gameClient.serverId]++

	return protocol.PackStart(startEvent), protocol.PackPlayer(startEvent.PlayerList[0])
}

// 실제 게임 오브젝트에 처리하는 로직
// Return이 1인 경우, BroadCast
func (gameClient *gameClient) handlePacket(packetData []byte) {

	var data []byte
	var msg []byte
	var sendType int
	packetType := packetData[0]

	switch packetType {
	case '2': // handleCommon(packetData[1:])
		packet := protocol.UnpackEvent(packetData[1:])
		gameClient.handleCommon(packet)
	case '4': // handleGameStart
		// 첫번째 데이터는 새로 입장하는 사람에게 전체 플레이어 리스트를 전달
		// 두번째 데이터는 입장한 유저에게 새로운 유저가 들어왔다는 데이터
		var player []byte
		data, player = gameClient.handleStart(packetData[1:])
		// 새로운 플레이어 데이터를 브로드케스팅
		msg = append([]byte{'3'}, player...)
		gameClient.broadcast <- msg
		// 방에 있는 유저들 리스트를 리스폰
		msg = append([]byte{'4'}, data...)
		sendType = 2
	case '5':
		data = gameClient.handleBullet(packetData[1:])
		msg = append([]byte{'5'}, data...)
		sendType = 1
	}

	// queue를 써야할까?
	switch sendType {
	case 1: gameClient.broadcast <- msg
	case 2: gameClient.eventGoing <- msg
	}

}
