package network

import (
	"net"
	"fmt"
	"github.com/GreyHood-Studio/play_server/utils"
	"github.com/GreyHood-Studio/play_server/model"
)

type GameServer struct {
	ServerId	int
	Port		int
	// 
	floor			model.Floor

	// 현재 접속 인원과 여태까지 접속했던 총 인원의 수
	maxConn		int
	currentConn	int

	// clients를 slice로 할까 map으로 할까 고민
	// map 가즈아
	clients map[string]*gameClient
	joins chan net.Conn
	broadcast chan []byte
}

func (server *GameServer) Broadcast(data []byte) {
	// loop문을 도는데 안도는게 베스트로 보임
	// if문 부하 방지
	if data[0] == 'q' {
		server.currentConn -= 1
		return
	}
	for _, client := range server.clients {
		client.outgoing <- data
	}
}

func (server *GameServer) ClientJoin(connection net.Conn) {
	// 새로운 클라이언트를 생성한 뒤, server의 clients list에 추가
	client := newClient(connection)
	fmt.Println("new client connection: ",client)
	// currentConn은 현재 붙어있는 Conn
	server.currentConn += 1

	// redis에서 토큰을 확인하여, 해당 id의 유저가 이 서버의 아이디인지 체크
	// 맨 처음 토큰에서의 data는 user의 ID
	data, err := client.reader.ReadBytes('\n')
	if err != nil && server.validateClient(string(data[:])) {
		fmt.Printf("Error client ID [%s]\n", string(data[:]))
		return
	}
	// 맞으면 추가
	fmt.Printf("Add new client ID [%s]\n", string(data[:]))
	server.clients[string(data[:])] = client
	client.Listen()

	go func() {
		for {
			server.broadcast <- <-client.broadcast
			}
	}()
}

func (server *GameServer) validateClient(clientID string) bool {

	return false
}

func (server *GameServer) handleRequest(data string) {
	// main process와 통신하는 IRC 채널 데이터 헨들링 함수
	// 여기서 select를 통해, 저 소켓을 닫을지 결정해야하는 것으로 보임 ( 작업 필요 )
	// select { } bool chan 하나 더 생성
	if data == "quit" {
		fmt.Printf("quit server port[%d]", server.Port)
		// 모든 클라이언트 종료 작업 ( graceful한 방식은 모든 서버들에게 데이터 죽는다는 전송 )
		quit := []byte{'q','u','i','t'}
		server.Broadcast(quit)
		// 바로 끄면 graceful하게 종료가 안되니, 클라이언트가 모두 종료되길 기다려야함
		for {
			if server.currentConn <= 0 {
				return
			}
		}
	}
}

func (server *GameServer) handleClient(packet Packet) Packet {

	return packet
}

func (server *GameServer) syncPacket() {
	// 주기적으로 계속 보내야하는 패킷에 대한 전달
	
}

func (server *GameServer) communicate() {
	go func() {
		for {
			for clientID, client := range server.clients {
				select {
				// 현재 map에 있는 list를 어케 routine을 돌게하지?
				case packet := <- client.packet:
					fmt.Println("communicate", clientID)
					client.packet <- server.handleClient(packet)
				}
			}
		}
	}()
}

func (server *GameServer) listen() {
	go func() {
		for {
			select {
			case data := <-server.broadcast:
				server.Broadcast(data)
			// 서버로부터 새로운 connection 요청이 들어오면, join 함수를 실행
			case conn := <-server.joins:
				server.ClientJoin(conn)
			}
		}
	}()
}

func (server *GameServer) Run() {
	// client들 listening
	server.listen()
	server.communicate()

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", server.Port))
	utils.CheckError(err, "Listen Socket Bind Error Check.")
	defer ln.Close() // main 함수가 끝나기 직전에 연결 대기를 닫음
	fmt.Printf("tcp server listen port: %d\n", server.Port)

	for {
		// Connection이 발생하면, connection을 server.joins 채널로 보냄
		conn, err := ln.Accept()
		utils.CheckError(err, "Accepted connection Error.")
		server.joins <- conn
	}
}

func (server *GameServer) Status() (int, int, int){
	// serverId, port, maxConn, currentConn
	return server.Port, server.maxConn, server.currentConn
}

func NewServer(serverId int, port int, maxConn int) *GameServer {
	// 서버 초기화
	server := &GameServer{
		ServerId: serverId,
		Port: port,
		maxConn: maxConn,

		currentConn: 0,

		clients: make(map[string]*gameClient),
		joins: make(chan net.Conn),
		broadcast: make(chan []byte),
	}

	return server
}