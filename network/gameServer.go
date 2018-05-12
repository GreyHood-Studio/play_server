package network

import (
	"net"
	"fmt"
	"github.com/GreyHood-Studio/server_util/error"
	"github.com/GreyHood-Studio/play_server/network/protocol"
	"github.com/GreyHood-Studio/play_server/model"
)

var floorMap map[int]*model.Floor

// 게임 서버가 곧 floor다. 맞지?
// client connection에 대한 정보만 전달하면 됨
// client는 gameServer의 id인 floor에 접근하면 됨
type GameServer struct {
	// 네트워크 로직과 관련 있는 곳
	ServerId	int
	InputPort	int
	EventPort	int
	maxConn		int

	// string은 clientName clients의 count는 len(server.clients)
	iJoins chan net.Conn
	eJoins chan net.Conn
	clients map[string]*gameClient
	tempClients map[string]*gameClient

	broadcast chan []byte
	inputcast chan []byte
	communicate chan []byte
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

func (server *GameServer) createClient(clientName string) *gameClient{
	// 새로운 클라이언트를 생성
	client := newClient(server.ServerId, clientName)

	client.clientName = clientName
	server.tempClients[clientName] = client

	return client
}

func (server *GameServer) addClient(clientName string, client *gameClient) {
	//server의 clients list에 추가
	client.clientId = len(server.clients) -1
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

	length, err := conn.Read(data)
	if length == 0 {
		fmt.Println("connected or disconnect socket in join")
		return
	}

	ctype, playerName, token := protocol.UnpackConnect(data)
	if err != nil && server.validateClient(playerName, token){
		fmt.Printf("Error Connect client Data [%s][%d]\n", length, string(data[:]))
		return
	}

	if client, ok := server.tempClients[playerName]; ok {
		// already exist logic
		//do something here
		if ctype == 0 {
			client.addEventConn(conn)
		} else if ctype == 1 {
			client.addInputConn(conn)
		}

		p := protocol.PackConnect(ctype, playerName)
		data = append([]byte{'1'}, p...)
		data = append(data, '\n')
		fmt.Printf("login data: %s", data)
		conn.Write(data)
		//conn.Write(data)
		server.addClient(playerName, client)

	} else {
		// playerName == clientName
		client = server.createClient(playerName)
		if ctype == 0 {
			client.addEventConn(conn)
		} else if ctype == 1 {
			client.addInputConn(conn)
		}
		p := protocol.PackConnect(ctype, playerName)
		data = append([]byte{'1'}, p...)
		data = append(data, '\n')
		fmt.Printf("login data: %s", data)
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
				fmt.Printf("broadcast data: %v\n%s\n", data, data)
				server.writeBroadcast(data)
			case data := <-server.inputcast:
				// router 역할만 수행
				fmt.Printf("inputcast data: %v\n%s\n", data, data)
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
	error.CheckError(err, "Listen Socket Bind Error Check.")
	defer eln.Close() // main 함수가 끝나기 직전에 연결 대기를 닫음
	iln, err := net.Listen("tcp", fmt.Sprintf(":%d", server.InputPort))
	error.CheckError(err, "Listen Socket Bind Error Check.")
	defer iln.Close() // main 함수가 끝나기 직전에 연결 대기를 닫음
	fmt.Printf("tcp server listen port: [%d][%d]\n", server.EventPort, server.InputPort)

	// event socket accpet
	go func() {
		for {
			// Connection이 발생하면, connection을 server.joins 채널로 보냄
			econn, err := eln.Accept()
			error.CheckError(err, "Accepted Event Socket connection Error.")
			server.eJoins <- econn
		}
	}()

	for {
		iconn, err := iln.Accept()
		error.CheckError(err, "Accepted Input Socket connection Error.")
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

func NewServer(serverId int, port int, maxConn int) *GameServer {
	// 서버 초기화
	if floorMap == nil {
		floorMap = make(map[int]*model.Floor)
	}

	floorMap[serverId] = model.NewFloor()

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
	}

	return server
}