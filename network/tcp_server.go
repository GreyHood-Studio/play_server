package network

import (
	"net"
	"fmt"
	"github.com/GreyHood-Studio/play_server/utils"
)

type TCPServer struct {
	ServerId	int
	Port		int
	maxConn		int

	// 현재 접속 인원과 여태까지 접속했던 총 인원의 수
	currentConn	int
	totalConn	int

	tcpChan		chan string

	clients []*Client
	joins chan net.Conn
	broadcast chan []byte
}

func (server *TCPServer) Broadcast(data []byte) {
	// loop문을 도는데 안도는게 베스트로 보임
	// if문 부하 방지
	if data[0] == 'e' && data[1] == 'x' {
		server.currentConn -= 1
	}
	for _, client := range server.clients {
		client.outgoing <- data
	}
}

func (server *TCPServer) ClientJoin(connection net.Conn) {
	// 새로운 클라이언트를 생성한 뒤, server의 clients list에 추가
	client := NewClient(connection, server.totalConn)
	fmt.Println("new client: ",client)
	// 증감 연산자가 안먹히네? 나중에 필요하면 변경
	// totalConn은 여태까지 접속했다 나갓다 들어온애들의 카운트
	// currentConn은 현재 붙어있는 Conn
	server.totalConn += 1
	server.currentConn += 1
	server.clients = append(server.clients, client)
	go func() {
		for { server.broadcast <- <-client.broadcast }
	}()
}

func (server *TCPServer) syncPacket() {
	// 주기적으로 계속 보내야하는 패킷에 대한 전달

}

func (server *TCPServer) handleRequest(data string) {
	// main process와 통신하는 IRC 채널 데이터 헨들링 함수
	// 여기서 select를 통해, 저 소켓을 닫을지 결정해야하는 것으로 보임 ( 작업 필요 )
	// select { } bool chan 하나 더 생성
	if data == "quit" {
		fmt.Println("quit server port[%d]", server.Port)
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

func (server *TCPServer) Listen() {
	go func() {
		for {
			select {
			case data := <-server.broadcast:
				server.Broadcast(data)
			// 서버로부터 새로운 connection 요청이 들어오면, join 함수를 실행
			case conn := <-server.joins:
				server.ClientJoin(conn)
			case data := <-server.tcpChan:
				server.handleRequest(data)
			}
		}
	}()
}

func (server *TCPServer) Run() {
	// client들 listening
	server.Listen()

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

func (server *TCPServer) Status() (int, int, int){
	// serverId, port, maxConn, currentConn
	return server.Port, server.maxConn, server.currentConn
}

// 서버 종료 로직 필요
func (server *TCPServer) Quit() {
	server.tcpChan <- "quit"
}

func NewServer(serverId int, port int, maxConn int) *TCPServer {
	// 서버 초기화
	server := &TCPServer{
		ServerId: serverId,
		Port: port,
		maxConn: maxConn,

		totalConn: 0,
		currentConn: 0,

		tcpChan: make(chan string),

		clients: make([]*Client, 0),
		joins: make(chan net.Conn),
		broadcast: make(chan []byte),
	}

	return server
}