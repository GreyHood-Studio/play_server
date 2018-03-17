package network


import (
	"bufio"
	"log"
	"net"
)

// Client holds info about connection
type TCP_Client struct {
	conn   net.Conn
	Server *tcp_server
}

// TCP server
type tcp_server struct {
	address                  string // Address to open connection: localhost:9999
	onNewClientCallback      func(c *TCP_Client)
	onClientConnectionClosed func(c *TCP_Client, err error)
	onNewMessage             func(c *TCP_Client, message string)
}

// Read client data from channel
func (c *TCP_Client) listen() {
	reader := bufio.NewReader(c.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			c.conn.Close()
			c.Server.onClientConnectionClosed(c, err)
			return
		}
		c.Server.onNewMessage(c, message)
	}
}

// Send text message to client
func (c *TCP_Client) Send(message string) error {
	_, err := c.conn.Write([]byte(message))
	return err
}

// Send bytes to client
func (c *TCP_Client) SendBytes(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

func (c *TCP_Client) Conn() net.Conn {
	return c.conn
}

func (c *TCP_Client) Close() error {
	return c.conn.Close()
}

// Called right after server starts listening new client
func (s *tcp_server) OnNewClient(callback func(c *TCP_Client)) {
	s.onNewClientCallback = callback
}

// Called right after connection closed
func (s *tcp_server) OnClientConnectionClosed(callback func(c *TCP_Client, err error)) {
	s.onClientConnectionClosed = callback
}

// Called when Client receives new message
func (s *tcp_server) OnNewMessage(callback func(c *TCP_Client, message string)) {
	s.onNewMessage = callback
}

// Start network server
func (s *tcp_server) Listen() {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal("Error starting TCP server.")
	}
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		client := &TCP_Client{
			conn:   conn,
			Server: s,
		}
		go client.listen()
		s.onNewClientCallback(client)
	}
}

// Creates new tcp server instance
func New(address string) *tcp_server {
	log.Println("Creating server with address", address)
	server := &tcp_server{
		address: address,
	}

	server.OnNewClient(func(c *TCP_Client) {})
	server.OnNewMessage(func(c *TCP_Client, message string) {})
	server.OnClientConnectionClosed(func(c *TCP_Client, err error) {})

	return server
}


