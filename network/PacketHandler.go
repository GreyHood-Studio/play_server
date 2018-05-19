package network

import (
	"fmt"
	"github.com/GreyHood-Studio/play_server/network/protocol"
	"github.com/GreyHood-Studio/play_server/model"
	"github.com/GreyHood-Studio/server_util/random"
)

func (gameClient *gameClient) handleCommon(commonBytes []byte) {
	event := protocol.UnpackEvent(commonBytes)

	// 후 정의
	switch event.EventType {

	}
}

func (gameClient *gameClient) handleHit(hitBytes []byte) []byte{
	hitPacket := protocol.UnpackHit(hitBytes)

	hitDamage := gameClient.room.Bullets[hitPacket.BulletNum].BulletDamage
	gameClient.room.Players[hitPacket.PlayerNum].CurrentHealth -=  hitDamage

	fmt.Printf("player[%d] damaged %d remain health %d\n",
		hitPacket.PlayerNum, hitDamage, gameClient.room.Players[hitPacket.PlayerNum].CurrentHealth)

	// player dead
	if gameClient.room.Players[hitPacket.PlayerNum].CurrentHealth <= 0 {
		hitPacket.EventType = 2
		//gameClient.exit()
	} else {
		hitPacket.EventType = 1
	}

	res := protocol.PackHit(hitPacket)

	return res
}

func (gameClient *gameClient) handleBullet(bulletBytes []byte) []byte{
	bulletNum := gameClient.room.LastBulletNum
	gameClient.room.LastBulletNum++
	bulletPacket := protocol.AssignBulletNum(bulletNum, bulletBytes)
	bulletDamage := protocol.ExtractBulletType(bulletBytes)

	gameClient.room.Bullets[bulletNum] = model.Bullet{BulletNum:bulletNum, BulletDamage:bulletDamage}
	return bulletPacket
}

// start packet은 완전히 나중에 변경할 protocol 데이터
func (gameClient *gameClient) handleStart(startBytes []byte) ([]byte, []byte){
	// +1을 하는 사람은 새로운 유저를 받을 수 있는 리스트를 만들기 위함
	var newPlayerNum int
	// 이건 동시성을 위해 가장 처음에 처리되어야 함
	newPlayerNum = gameClient.room.LastPlayerNum
	gameClient.room.LastPlayerNum++
	// deep copy logic
	existPlayers := make(map[int]*model.Player)
	// deep copy
	for key, value := range gameClient.room.Players {
		existPlayers[key] = value
	}

	// 새로운 유저 추가
	newPlayer := protocol.ExtractPlayerInStartPacket(startBytes)
	newPlayer.PlayerNum = newPlayerNum
	gameClient.room.Players[newPlayerNum] = newPlayer
	// network game client에 client id 부여
	gameClient.clientId = newPlayerNum

	spawnId := random.CreateRandomInteger(1)

	responsePacket := protocol.NewStartPacket(gameClient.room.MapId, spawnId, 1)
	responsePacket.ConnectionInfo[0] = newPlayer

	// player packet에서 항목 추출
	// 새로 추가된 player까지 추가 됨
	for _, value := range existPlayers {
		// 기존 유저의 리스트를 받아옴
		fmt.Printf("current user[%d]: %s\n",value.PlayerNum, value.PlayerName)
		responsePacket.ConnectionInfo = append(responsePacket.ConnectionInfo, value)
	}

	if gameClient.room.LastPlayerNum -1 != newPlayerNum {
		fmt.Printf("access same time different player id %d %d\n", gameClient.room.LastPlayerNum, newPlayerNum)
	}

	res := protocol.PackStart(*responsePacket)
	// 처음꺼만 남김
	responsePacket.ConnectionInfo = responsePacket.ConnectionInfo[:1]
	broad := protocol.PackStart(*responsePacket)

	return res, broad
}

func (gameClient *gameClient) handleMoveRoom(packetData []byte) {
	// 방 이동 로직

	// redis cache를 확인하여, aggregation으로 유저에게 방을 추천
	// 해당 방 어드레스 전달

}

// 실제 게임 오브젝트에 처리하는 로직
// Return이 1인 경우, BroadCast
func (gameClient *gameClient) handlePacket(packetData []byte) {

	var data []byte
	var msg []byte
	var sendType int
	packetType := packetData[0]

	switch packetType {
	case '1':
		// 입력 이벤트
	case '2':
		// 총알 발사 이벤트
		data = gameClient.handleBullet(packetData[1:])
		msg = append([]byte{'2'}, data...)
		sendType = 1
	case '3':
		// 총알 피격 이벤트
		data = gameClient.handleHit(packetData[1:])
		msg = append([]byte{'3'}, data...)
		sendType = 1
	case '4': // handleGameStart
		// 첫번째 데이터는 새로 입장하는 사람에게 전체 플레이어 리스트를 전달
		// 두번째 데이터는 입장한 유저에게 새로운 유저가 들어왔다는 데이터
		data, broad := gameClient.handleStart(packetData[1:])
		// 새로운 플레이어 데이터를 브로드케스팅
		broad = append([]byte{'4'}, broad...)
		gameClient.broadcast <- broad
		// 방에 있는 유저들 리스트를 리스폰
		msg = append([]byte{'4'}, data...)
		sendType = 2
	case '5':
		// 공통 패킷
		gameClient.handleCommon(packetData[1:])
	case '9':
		// 방 이동 패킷
		gameClient.handleMoveRoom(packetData[1:])
	}

	// queue를 써야할까?
	switch sendType {
	case 1: gameClient.broadcast <- msg
	case 2: gameClient.eventGoing <- msg
	}
}