package main

import (
	"go.uber.org/zap"
	"github.com/GreyHood-Studio/play_server/network"
	"github.com/gin-gonic/gin"
)

// game server handler
func mainTCPServer(port string) {
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

func mainHTTPServer(port string) {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	r.Run(port)
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	logger.Info("failed to fetch URL",
		// Structured context as strongly typed Field values.
		zap.String("test", "testdata"),
	)
	confPort := readDefaultConfig()

	mainTCPServer(confPort)
	mainHTTPServer(":5000")
}