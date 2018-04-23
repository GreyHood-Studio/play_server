package network

// 실제 정의된 protocol을 처리하는 구역
func (server *GameServer) handlePacket(packet Packet) Packet {
	// case 1: game request

	switch packet.MsgFormat {
	case 1: result := server.floor.GetFloorStatus()
		packet.MsgBody = result
		packet.MsgFormat = 2
		packet.MsgFormat = 3
	}

	return packet
}

// 실제 게임 오브젝트에 처리하는 로직
// Return이 1인 경우, BroadCast
func (gameClient *gameClient) handlePacket(packetByte []byte, sendType int) {

	switch sendType {
	// -3: fatalError
	// -2: Error
	// -1: No Receive error
	// 0: broadcast
	// 1: request
	// 2: response
	// 3: response&broadcast
	case -3: gameClient.exit()
	case -2: gameClient.outgoing <- []byte{'q','\n'}
	case 0: gameClient.broadcast <- packetByte
	case 1: gameClient.outgoing <- packetByte
	}
}