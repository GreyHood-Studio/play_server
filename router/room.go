package router

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/GreyHood-Studio/play_server/network"
	"net/http"
	"fmt"
)

// 전체 floorCount
var maxRoomCount		int
var roomMap			map[int]network.GameServer

// goroutine으로 돌면서 floor들 데이터를 지속적으로 관리해주는 함수
func ManageFoom(maxCount int) {
	maxRoomCount =	maxCount
	roomMap = make(map[int]network.GameServer)
}

// Floor의 정보를 가져오는 로직
// 만약에 Redis에 있으면 필요 없음
func getRoom(c *gin.Context) {
	ServerIDString := c.Param("serverID")
	serverId, _ := strconv.Atoi(ServerIDString)

	// map 키 체크
	_, exists := roomMap[serverId]
	if exists {

	} else {
		c.String(http.StatusNotFound, fmt.Sprintf("Don't Exist ServerId %d", serverId))
		return
	}
}

// 입력데이터는 업데이트만 하는거고 주기적으로 update
// Floor를 생성하는 로직 -> 실제 게임 서버 소켓 생성 로직
func createRoom(c *gin.Context) {
	// 실제 게임 서버의 층 오픈
	// Packet 층 말고도, Floor의 정보를 보관하는 로직이 필요
	ServerIDString := c.Param("serverID")
	portString := c.PostForm("port")
	maxConnString := c.PostForm("maxConn")
	port, pErr := strconv.Atoi(portString)
	serverId, sErr := strconv.Atoi(ServerIDString)
	maxConn, mErr := strconv.Atoi(maxConnString)

	if pErr != nil || sErr != nil || mErr != nil {
		c.String(http.StatusBadRequest,
			fmt.Sprintf("Invalid Parameter create Floor Server[%d] Port[%d]", serverId, port))
		return
	}
	fmt.Println(port, serverId, maxConn)

	// server create validation logic
	if len(roomMap) >= maxRoomCount {
		c.String(http.StatusNotAcceptable, fmt.Sprintf("Already Max Server Count"))
		return
	}

	_, exists := roomMap[serverId]
	//if 포트 또는 서버아이디가 중복된 경우 Return
	if exists {
		c.String(http.StatusNotAcceptable,
			fmt.Sprintf("Already Exist Server Port %d OR ServerId %d", port, serverId))
		return
	}
	// server validation logic complete

	server	:= network.NewServer(serverId, port, maxConn)
	go server.Run()
	// 서버 map에 serverID를 기준으로 서버를 생성
	roomMap[serverId] = *server

	// 현재 생성된 서버의 수
	fmt.Printf("currentFloorCount %d\n", len(roomMap))
	c.String(http.StatusOK,
		fmt.Sprintf("Success Create Floor[%d] Server[%d] Port[%d]", len(roomMap), serverId, port))
}

// floor를 삭제하는 로직 ( Refresh 또는 맵 구조 변경용?
func deleteRoom(c *gin.Context) {
	ServerIDString := c.Param("serverID")
	serverId, err := strconv.Atoi(ServerIDString)
	if err != nil {
		c.String(http.StatusBadRequest,
			fmt.Sprintf("Invalid Parameter Delete Floor ServerId %d", serverId))
		return
	}
	delete(roomMap, serverId)
	c.String(http.StatusOK, fmt.Sprintf("Delete Floor Floor %d",  serverId))
}