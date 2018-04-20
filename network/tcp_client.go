package network

import (
	"bufio"
	"github.com/GreyHood-Studio/play_server/utils"
	"github.com/GreyHood-Studio/play_server/network/protocol"
	"net"
)

type Client struct {
	// clientID	int	 ClientID가 필요할까? 필요한 경우 삽입
	// host/port/serverId/ClientId로 Client들 식별 가능
	// 나중에는 client와 mapping 필요
	clientId 	int
	conn		net.Conn
	broadcast 	chan []byte
	outgoing 	chan []byte
	reader   	*bufio.Reader
	writer   	*bufio.Writer
}

func (client *Client) Read() {
	for {
		data, err := client.reader.ReadBytes('\n')
		utils.CheckError(err, "Client disconnected.\n")
		//route_packet에서 모든 로직을 처리한 뒤에 반환
		// client 종료 알람
		if data[0] == 'q' && data[1] =='u'{
			client.exit()
		}

		result, sendType := protocol.RoutePacket(data)

		switch sendType {
		// 0: broadcast, 1: return 2: quit
		case 0: client.broadcast <- result
		case 1: client.outgoing <- result
		}
	}
}

func (client *Client) Write() {
	for data := range client.outgoing {
		client.writer.Write(data)
		client.writer.Flush()
	}
}

func (client *Client) Listen() {
	go client.Read()
	go client.Write()
}

func (client *Client) exit() {
	quit := []byte{'e','x','i','t'}
	client.broadcast <- quit
	client.conn.Close()
}

func NewClient(connection net.Conn, clientId int) *Client {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	client := &Client{
		conn: connection,
		clientId: clientId,
		broadcast: make(chan []byte),
		outgoing: make(chan []byte),
		reader: reader,
		writer: writer,
	}

	client.Listen()

	return client
}