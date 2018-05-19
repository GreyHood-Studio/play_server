package main

import (
	"github.com/GreyHood-Studio/play_server/router"
	"github.com/GreyHood-Studio/server_util/setup"
	"fmt"
)

func main() {
	config := setup.NewConfig()
	fmt.Println(config)

	setup.ConnectDatabase(config.Database)
	r := setup.ConnectRouter(config.Cache)

	router.SetAPIRoute(r,3)
	// 게임 서버에서 Web Server와 통신하기 위한 HTTP Server Port Open
	r.Run(fmt.Sprint(":", config.Server.Port))
}