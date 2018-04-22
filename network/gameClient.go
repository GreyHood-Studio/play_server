package network

import (
	"github.com/GreyHood-Studio/play_server/utils"
	"bufio"
	"net"
	"fmt"
	"reflect"
)

type gameClient struct {
	// gameClientID	int	 gameClientID가 필요할까? 필요한 경우 삽입
	// host/port/serverId/gameClientId로 gameClient들 식별 가능
	// 나중에는 gameClient와 mapping 필요
	clientName	string
	clientId 	int
	conn		net.Conn

	// 정말 packet을 주고받을지 고민
	packet 		chan Packet

	broadcast 	chan []byte
	outgoing 	chan []byte
	reader   	*bufio.Reader
	writer   	*bufio.Writer
}

// 실제 게임 오브젝트에 처리하는 로직
// Return이 1인 경우, BroadCast
func (gameClient *gameClient) handlePacket(packet Packet) {
	switch packet.MsgFormat {
	case 0: fmt.Println("ping check")
	case 1:
		//return protocol.RequestGameStart(packet.MsgBody), 2
	case 2:
	case 3:
	}

	//return []byte{'4'}, 0
}

// 에러가 아닌 패킷을 처리하는 로직
func (gameClient *gameClient) handleClient(data []byte) ([]byte, int) {
	// Json으로 패킷을 Unmarshal 하는 로직
	packet := unpackPacket(data)
	if packet.MsgType == -2 {
		// unmarshal error
		return []byte{'e','r','r','o','r','\n'}, -1
	}

	// packet을 모두 처리한 뒤에, receive가 필요한 경우, marshaling을 하는 로직
	result, sendType := packRawPacket(packet)
	
	return result, sendType
}

func (gameClient *gameClient) Read() {
	for {
		data, err := gameClient.reader.ReadBytes('\n')
		if utils.NoDeadError(err, "gameClient disconnected.\n") {
			gameClient.exit()
			return
		}

		// 디버깅용 클라이언트 데이터 체크
		fmt.Printf("gameClient[%d] type %v msg[0] type %v msg : %s",
			gameClient.clientId, reflect.TypeOf(data), reflect.TypeOf(data[0]), data)

		//route_packet에서 모든 로직을 처리한 뒤에 반환
		// gameClient 종료 알람
		if data[0] == 'q' {
			gameClient.exit()
		}

		result, sendType := gameClient.handleClient(data)

		switch sendType {
		// -3: fatalError -2: Error -1: No Receive error 0: broadcast, 1: request 2: response, 3: response&broadcast
		case -3: gameClient.exit()
		case -2: gameClient.outgoing <- []byte{'q','\n'}
		case -1: return
		//case -1: gameClient.outgoing <- data
		case 0: gameClient.broadcast <- result
		case 1: gameClient.outgoing <- result
		}
	}
}

func (gameClient *gameClient) Write() {
	for data := range gameClient.outgoing {
		gameClient.writer.Write(data)
		gameClient.writer.Flush()
	}
}

func (gameClient *gameClient) Listen() {
	go gameClient.Read()
	go gameClient.Write()
}

func (gameClient *gameClient) exit() {
	quit := []byte{'q','u','i','t'}
	gameClient.broadcast <- quit
	gameClient.conn.Close()
}

func newClient(connection net.Conn) *gameClient {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	gameClient := &gameClient{
		conn: connection,
		broadcast: make(chan []byte),
		outgoing: make(chan []byte),
		packet: make(chan Packet),
		reader: reader,
		writer: writer,
	}

	return gameClient
}