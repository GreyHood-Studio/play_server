package generator

import "github.com/GreyHood-Studio/play_server/network"

func CreateFloor(floorId string) {

}

// game server handler
func bindFloorPort(port string) {
	server := network.New(port)

	server.OnNewClient(func(c *network.TCP_Client) {
		// new client connected
		// lets send some message
		c.Send("Hello")
	})
	server.OnNewMessage(func(c *network.TCP_Client, message string) {
		// new message received
	})
	server.OnClientConnectionClosed(func(c *network.TCP_Client, err error) {
		// connection
		// with client lost
	})

	server.Listen()
}
