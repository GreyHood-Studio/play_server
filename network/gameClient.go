package network

import (
	"bufio"
	"net"
	"fmt"
	"reflect"
	"github.com/GreyHood-Studio/play_server/network/protocol"
	"github.com/GreyHood-Studio/server_util/error"
)

type gameClient struct {
	// gameClientID	int	 gameClientID가 필요할까? 필요한 경우 삽입
	// host/port/serverId/gameClientId로 gameClient들 식별 가능
	// 나중에는 gameClient와 mapping 필요
	clientName	string
	clientId 	int
	serverId	int
	inputConn	net.Conn
	eventConn	net.Conn

	// ipc를 위한 채널
	communicate	chan []byte
	broadcast 	chan []byte
	inputcast	chan []byte
	inputGoing 	chan []byte
	eventGoing 	chan []byte
	// socket read write에 대한 처리
	buf			[]byte
	iReader		*bufio.Reader
	eReader   	*bufio.Reader
}


func (gameClient *gameClient) inputRead() {
	for {
		data, err := gameClient.iReader.ReadBytes('\n')
		if error.NoDeadError(err, "gameClient disconnected.\n") {
			gameClient.exit()
			return
		}
		gameClient.inputcast <- data
	}
}

func (gameClient *gameClient) inputWrite() {
	for data := range gameClient.inputGoing {
		gameClient.inputConn.Write(data)
		//fmt.Printf("input write data %s",data)
	}
}

func (gameClient *gameClient) eventRead() {
	for {
		data, err := gameClient.eReader.ReadBytes('\n')
		if error.NoDeadError(err, "gameClient disconnected.\n") {
			gameClient.exit()
			return
		}

		// 디버깅용 클라이언트 데이터 체크
		fmt.Printf("gameClient[%d] type_%v: msg[0]_type: %v msg: %s",
			gameClient.clientId, reflect.TypeOf(data), reflect.TypeOf(data[0]), data)

		go gameClient.handlePacket(data)
	}
}

func (gameClient *gameClient) eventWrite() {
	for data := range gameClient.eventGoing {
		gameClient.buf = append(data, '\n')
		gameClient.eventConn.Write(gameClient.buf)
		fmt.Printf("eventWrite [%d] %s",gameClient.clientId, gameClient.buf)
		gameClient.buf = nil
	}
}

func (gameClient *gameClient) Listen() {
	// input read와 event read를 구분할 것
	go gameClient.inputRead()
	go gameClient.inputWrite()
	go gameClient.eventRead()
	go gameClient.eventWrite()
}

func (gameClient *gameClient) exit() {
	// exit packet
	data := protocol.PackEvent(6, gameClient.clientId, 0)

	println("quit client ", gameClient.clientName)
	roomMap[gameClient.serverId].DeletePlayer(gameClient.clientId)

	gameClient.inputConn.Close()
	gameClient.eventConn.Close()

	gameClient.communicate <- append([]byte{'q'}, []byte(gameClient.clientName)...)
	gameClient.broadcast <- append([]byte{'2'}, []byte(data)...)
}

func (gameClient *gameClient)addEventConn(conn net.Conn) {
	reader := bufio.NewReader(conn)

	gameClient.eReader = reader
	gameClient.eventConn = conn
	fmt.Println("addEventConn in client")
}

func (gameClient *gameClient)addInputConn(conn net.Conn) {
	reader := bufio.NewReader(conn)

	gameClient.iReader = reader
	gameClient.inputConn = conn
	fmt.Println("addInputConn in client")
}

func newClient(serverId int, clientName string) *gameClient {

	gameClient := &gameClient{
		clientName: clientName,
		serverId: serverId,

		buf: make([]byte, 2048),

		broadcast: make(chan []byte),
		inputcast: make(chan []byte),
		eventGoing: make(chan []byte),
		inputGoing: make(chan []byte),
		// always broadcast in game server
	}

	return gameClient
}