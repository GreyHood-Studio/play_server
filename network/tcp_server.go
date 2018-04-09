package network

import (
	"net"
	"fmt"
	"bufio"
)

type TCPServer struct {
	FloorNum	int
	Port		int
	conns 		[]net.Conn
	channels	chan ClientJob
}

type ClientJob struct {
	data []byte
	conn net.Conn
}

// 모든 클라이언트한테 전달해주는 데이터
func (tcp *TCPServer) notifyClient() {

}

// 특정 유저한테 전달해주는 데이터
func (tcp *TCPServer) commandClient() {

}

// 유저의 Request 없이도 주기적으로 돌아서 처리해주는 로직
func (tcp *TCPServer) updateClient(conn net.Conn) {
	// 처리해야하는 로직

}

func (tcp *TCPServer) responseClient() {
	for {
		// Wait for the next job to come off the queue.
		clientJob := <-tcp.channels

		// CPU가 처리하는 작업
		// Do something thats keeps the CPU buys for a whole second.
		parseData(clientJob.data)

		// 데이터 보낸 사람에게 응답 보낼 경우,
		clientJob.conn.Write([]byte("ddd"))
		// 특정 유저한테 응답 보낼 경우
		// ????
	}
}

func (tcp *TCPServer) requestClient(conn net.Conn) {
	buf := bufio.NewReader(conn)

	for {
		// data
		data, err := buf.ReadBytes('\n')
		if err != nil {
			fmt.Printf("Client disconnected.\n")
			break
		}

		tcp.channels <- ClientJob{data, conn}
	}
}

func (tcp *TCPServer) Run() {
	// Client Job
	tcp.channels = make(chan ClientJob)
	// request의 response를 처리하는 socket list
	go tcp.responseClient()

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", tcp.port))
	check(err, "Accepted connection.")

	defer ln.Close() // main 함수가 끝나기 직전에 연결 대기를 닫음

	for {
		conn, err := ln.Accept() // 클라이언트가 연결되면 TCP 연결을 리턴
		check(err, "Accepted connection.")

		// Request를 처리하기 위한 Socket List
		go tcp.requestClient(conn)
		go tcp.updateClient(conn)
	}
}