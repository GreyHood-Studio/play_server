package route

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/GreyHood-Studio/play_server/network"
	"net/http"
	"fmt"
	"github.com/GreyHood-Studio/play_server/utils"
)

type serverInfo struct {
	serverId	int
	port		int
	server		network.TCPServer
}

// goroutine으로 돌면서 floor들 데이터를 지속적으로 관리해주는 함수
func ManageFloor() {
	//servers := make([]*serverInfo, 0)
	//
	//for {
	//
	//}
}

// Floor의 정보를 가져오는 로직
// 만약에 Redis에 있으면 필요 없음
func getFloorStatus(c *gin.Context) {
	//ServerIDString := c.Param("serverId")
}

// 입력데이터는 업데이트만 하는거고 주기적으로 update
// Floor를 생성하는 로직 -> 실제 게임 서버 소켓 생성 로직
func createFloor(c *gin.Context) {
	// 실제 게임 서버의 층 오픈
	ServerIDString := c.Param("serverId")
	portString := c.PostForm("port")
	maxConnString := c.PostForm("maxConn")
	port, err := strconv.Atoi(portString)
	utils.CheckError(err, "port arg error")
	serverId, err := strconv.Atoi(ServerIDString)
	utils.CheckError(err, "serverId arg error")
	maxConn, err := strconv.Atoi(maxConnString)
	utils.CheckError(err, "maxConnection arg error")
	//mapType := c.Param("mapType")
	//maxUser := c.Param("maxUser")

	// floorServer List를 보관하는 로직이 필요 -> Redis로 전달해야할꺼같은데?
	server := network.NewServer(port, serverId, maxConn)
	go server.Run()

	c.String(http.StatusOK, fmt.Sprintf("Create FloorServer Port %d Floor %d", port, serverId))
}

// floor를 삭제하는 로직 ( Refresh 또는 맵 구조 변경용?
func deleteFloor(c *gin.Context) {
	ServerIDString := c.Param("serverId")
	c.String(http.StatusOK, fmt.Sprintf("Delete FloorServer Floor %s",  ServerIDString))

	//server = nil
}