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


func (gameClient *gameClient) handleClient(packet Packet) {
	requestByte, sendType := packRawPacket(packet)

	gameClient.handlePacket(requestByte, sendType)
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

		gameClient.handleClient(unpackPacket(data))
	}
}

func (gameClient *gameClient) Write() {
	for data := range gameClient.outgoing {
		gameClient.writer.Write(data)
		gameClient.writer.Flush()
	}
}

func (gameClient *gameClient) communicate() {
	go func() {
		for {
			select {
				case packet := <- gameClient.packet:
					gameClient.handleClient(packet)
			}
		}
	}()
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