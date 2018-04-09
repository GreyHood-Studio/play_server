package route

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/GreyHood-Studio/play_server/network"
)

// Floor의 정보를 가져오는 로직
// 만약에 Redis에 있으면 필요 없음
func getFloorStatus(c *gin.Context) {

}

// 입력데이터는 업데이트만 하는거고 주기적으로 update

// Floor를 생성하는 로직 -> 실제 게임 서버 소켓 생성 로직
func createFloor(c *gin.Context) {
	// 실제 게임 서버의 층 오픈
	portString := c.Param("port")
	port, _ := strconv.Atoi(portString)
	//mapType := c.Param("mapType")
	//maxUser := c.Param("maxUser")

	

	// floorServer List를 보관하는 로직이 필요 -> Redis로 전달해야할꺼같은데?
	floorServer := network.TCPServer{Port:port,FloorNum:1}
	floorServer.Run()
}

// floor를 삭제하는 로직 ( Refresh 또는 맵 구조 변경용?
func deleteFloor(c *gin.Context) {

}