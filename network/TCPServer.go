package network

import (
	"net"
	"fmt"
	"github.com/GreyHood-Studio/server_util/checker"
	"github.com/GreyHood-Studio/play_server/network/protocol"
	"github.com/GreyHood-Studio/play_server/model"
	"errors"
	"bufio"
)

// 게임 서버의 네트워크 정보
type GameServer struct {
	// 네트워크 로직과 관련 있는 곳
	ServerId	int
	InputPort	int
	EventPort	int
	maxConn		int

	// string은 clientName clients의 count는 len(server.clients)
	iJoins chan net.Conn
	eJoins chan net.Conn
	// game client들의 리스트 string은 player_name
	clients map[string]*gameClient
	// 완전히 접속하기 직전의 임시 클라이언트
	tempClients map[string]*gameClient

	broadcast chan []byte
	inputcast chan []byte
	communicate chan []byte

	Room *model.Room
}

func (server *GameServer) writeBroadcast(data []byte) {
	// loop문을 도는데 안도는게 베스트로 보임
	for _, client := range server.clients {
		client.eventGoing <- data
	}
}

func (server *GameServer) writeInputcast(data []byte) {
	// loop문을 도는데 안도는게 베스트로 보임
	for _, client := range server.clients {
		client.inputGoing <- data
	}
}

// Create Client 생성
func (server *GameServer) createClient(clientName string) *gameClient{
	// 새로운 클라이언트를 생성
	client := newClient(server.ServerId, clientName, server.Room)

	client.clientName = clientName
	server.tempClients[clientName] = client

	return client
}

func (server *GameServer) addClient(clientName string, client *gameClient) {
	//server의 clients list에 추가
	client.communicate = server.communicate
	server.clients[clientName] = client

	fmt.Printf("m[%d]/c[%d] Add new client ID [%s]\n", server.maxConn, len(server.clients), clientName)
	client.Listen()

	go func() {
		for {
			server.broadcast <- <-client.broadcast
		}
	}()

	go func() {
		for {
			server.inputcast <- <- client.inputcast
		}
	}()
}

func (server *GameServer) clientJoin(conn net.Conn) {
	data := make([]byte, 1024)

	reader := bufio.NewReader(conn)
	data, err := reader.ReadBytes('\n')
	if checker.NoDeadError(err, "gameClient disconnected.\n") {
		conn.Close()
	}
	//length, err := conn.Read(data)
	//if length == 0 {
	//	fmt.Println("connected or disconnect socket in join")
	//	return
	//}
	if data[0] != '0' {
		checker.NoDeadError(errors.New("client_init_packet_format_error"), string(data))
	}

	initPacket := protocol.UnpackInitPacket(data[1:])
	if err != nil && server.validateClient(initPacket.PlayerName, initPacket.Token){
		fmt.Printf("Error Connect client Data [%s]\n", string(data[:]))
		return
	}

	if client, ok := server.tempClients[initPacket.PlayerName]; ok {
		// already exist logic
		//do something here
		if initPacket.ConnectType == 0 {
			client.addEventConn(conn)
		} else if initPacket.ConnectType == 1 {
			client.addInputConn(conn)
		}

		p := protocol.PackInitPacket(initPacket.ConnectType, initPacket.PlayerName)
		data = append([]byte{'0'}, p...)
		data = append(data, '\n')
		conn.Write(data)
		//conn.Write(data)
		server.addClient(initPacket.PlayerName, client)

		// Redis에 클라이언트 추가된 내용을 삽입
	} else {
		// playerName == clientName
		client = server.createClient(initPacket.PlayerName)
		if initPacket.ConnectType == 0 {
			client.addEventConn(conn)
		} else if initPacket.ConnectType == 1 {
			client.addInputConn(conn)
		}
		p := protocol.PackInitPacket(initPacket.ConnectType, initPacket.PlayerName)
		data = append([]byte{'0'}, p...)
		data = append(data, '\n')
		conn.Write(data)
		//conn.Write(data)
	}
}

func (server *GameServer) validateClient(playerName string, token string) bool {
	// redis에서 토큰을 확인하여, 해당 id의 유저가 이 서버의 아이디인지 체크
	return false
}

func (server *GameServer) handleBroadcast() {
	go func() {
		for {
			select {
			// client 로부터 broadcast 데이터를 받는 곳
			case data := <-server.communicate:
				if data[0] == 'q' {
					server.closeClient(string(data[1:]))
				}
			case data := <-server.broadcast:
				server.writeBroadcast(data)
			case data := <-server.inputcast:
				// router 역할만 수행
				//fmt.Printf("inputcast data: %v\n%s", data, data)
				server.writeInputcast(data)
			}
		}
	}()
}

func (server *GameServer) listen() {
	go func() {
		for {
			select {
			case conn := <-server.eJoins:
				// 서버로부터 새로운 connection 요청이 들어오면, join 함수를 실행
				fmt.Println("connect event socket")
				server.clientJoin(conn)
			case conn := <-server.iJoins:
				fmt.Println("connect input socket")
				// 서버로부터 새로운 connection 요청이 들어오면, join 함수를 실행
				server.clientJoin(conn)
			}
		}
	}()
}

func (server *GameServer) Run() {
	// client들 listening
	server.handleBroadcast()
	server.listen()

	eln, err := net.Listen("tcp", fmt.Sprintf(":%d", server.EventPort))
	checker.CheckError(err, "Listen Socket Bind Error Check.")
	defer eln.Close() // main 함수가 끝나기 직전에 연결 대기를 닫음
	iln, err := net.Listen("tcp", fmt.Sprintf(":%d", server.InputPort))
	checker.CheckError(err, "Listen Socket Bind Error Check.")
	defer iln.Close() // main 함수가 끝나기 직전에 연결 대기를 닫음
	fmt.Printf("tcp server listen port: [%d][%d]\n", server.EventPort, server.InputPort)

	// event socket accpet
	go func() {
		for {
			// Connection이 발생하면, connection을 server.joins 채널로 보냄
			econn, err := eln.Accept()
			checker.CheckError(err, "Accepted Event Socket connection Error.")
			server.eJoins <- econn
		}
	}()

	for {
		iconn, err := iln.Accept()
		checker.CheckError(err, "Accepted Input Socket connection Error.")
		server.iJoins <- iconn
	}
}

func (server *GameServer) closeClient(clientName string)  {
	delete(server.clients, clientName)
	delete(server.tempClients, clientName)
}

func (server *GameServer) Status() (int, int, int, int){
	// serverId, port, maxConn, currentConn
	return server.InputPort, server.EventPort, server.maxConn, len(server.clients)
}

func NewServer(serverId int, port int, maxConn int, room *model.Room) *GameServer {
	// 서버 초기화
	server := &GameServer{
		ServerId: serverId,
		EventPort: port,
		InputPort: port+1,

		maxConn: maxConn,

		iJoins: make(chan net.Conn),
		eJoins: make(chan net.Conn),
		clients: make(map[string]*gameClient),
		tempClients: make(map[string]*gameClient),

		broadcast: make(chan []byte),
		inputcast: make(chan []byte),
		communicate: make(chan []byte),

		Room: room,
	}

	return server
}