package main

import (
	"go.uber.org/zap"
	"github.com/gin-gonic/gin"
	"github.com/GreyHood-Studio/play_server/route"
)

func mainHTTPServer(port string) {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	route.SetAPIRoute(r)

	r.Run(port)
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	logger.Info("failed to fetch URL",
		// Structured context as strongly typed Field values.
		zap.String("test", "testdata"),
	)
	//confPort := readDefaultConfig()

	// 게임 서버에서 Web Server와 통신하기 위한 HTTP Server Port Open
	mainHTTPServer(":5000")
}