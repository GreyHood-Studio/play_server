package network

import (
	"github.com/GreyHood-Studio/play_server/network/protocol"
)


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

func (gameClient *gameClient) handlePlayer(event protocol.PlayerEvent) {

}

func (gameClient *gameClient) handleStart(packet protocol.StartEvent) []byte{
	players := floorMap[gameClient.serverId].Players
	for pid, value := range players {
		playerEvent := protocol.PlayerEvent{ PlayerName: value.PlayerName, PlayerId:pid }
		packet.PlayerList = append(packet.PlayerList, playerEvent)
	}
	floorMap[gameClient.serverId].AddPlayer(gameClient.clientName, gameClient.clientId)
	packet.PlayerList[0].PlayerId = len(packet.PlayerList)
	packet.PlayerList[0].PlayerName = gameClient.clientName

	return protocol.PackStart(packet)
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
	case '3': // handleNewPlayer
		packet := protocol.UnpackPlayer(packetData[1:])
		gameClient.handlePlayer(packet)
	case '4': // handleGameStart
		packet := protocol.UnpackStart(packetData[1:])
		sendType = 1
		data = gameClient.handleStart(packet)
		msg = append([]byte{'4'}, data...)
	}

	// queue를 써야할까?
	switch sendType {
	case 1: gameClient.broadcast <- msg
	case 2: gameClient.eventGoing <- msg
	}

}
