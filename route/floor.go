package route

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/GreyHood-Studio/play_server/network"
	"github.com/GreyHood-Studio/play_server/model"
	"net/http"
	"fmt"
)

// 전체 floorCount
var currentFloorCount	int
var maxFloorCount		int
var serverMap			map[int]gameServer

type gameServer struct {
	tcpServer		network.TCPServer
	floor			model.Floor

	tcpChan			chan string
	floorChan		chan string
	// mainserver에서 gameserver들에게 주는 요청
}

// goroutine으로 돌면서 floor들 데이터를 지속적으로 관리해주는 함수
func ManageFloor(maxCount int) {
	currentFloorCount = 0
	maxFloorCount =	maxCount
	//servers = make([]network.TCPServer, maxFloorCount)
	serverMap = make(map[int]gameServer)

	for {
		select {
	//	case data := floors[0].Share:
	//		fmt.Println(data)
		}
	}
}

// Floor의 정보를 가져오는 로직
// 만약에 Redis에 있으면 필요 없음
func getFloor(c *gin.Context) {
	ServerIDString := c.Param("serverId")
	serverId, _ := strconv.Atoi(ServerIDString)

	// map 키 체크
	_, exists := serverMap[serverId]
	if exists {
		//c.String(http.StatusOK, fmt.Sprintf("server[%d] map[%d] currentUser %d",id,
		//	server.floor.Status.MapType, server.floor.Status.PlayerCount))
	} else {
		c.String(http.StatusNotFound, fmt.Sprintf("Don't Exist ServerId %d", serverId))
		return
	}
}

// 입력데이터는 업데이트만 하는거고 주기적으로 update
// Floor를 생성하는 로직 -> 실제 게임 서버 소켓 생성 로직
func createFloor(c *gin.Context) {
	// 실제 게임 서버의 층 오픈
	// Packet 층 말고도, Floor의 정보를 보관하는 로직이 필요
	ServerIDString := c.Param("serverId")
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

	if currentFloorCount >= maxFloorCount {
		c.String(http.StatusNotAcceptable, fmt.Sprintf("Already Max Server Count"))
		return
	}

	_, exists := serverMap[serverId]
	//if 포트 또는 서버아이디가 중복된 경우 Return
	if exists {
		c.String(http.StatusNotAcceptable,
			fmt.Sprintf("Already Exist Server Port %d OR ServerId %d", port, serverId))
		return
	}
	currentFloorCount += 1

	// 필요할까? 고민좀 해봅시다
	tcpChan	:= make(chan string)
	floorChan := make(chan string)

	// floorServer List를 보관하는 로직이 필요 -> Redis로 전달해야할꺼같은데?
	server	:= network.NewServer(serverId, port, maxConn)
	floor	:= model.Floor{}

	go server.Run()
	go floor.Run()

	serverMap[serverId] = gameServer{*server, floor, tcpChan, floorChan}

	fmt.Println("currentFloorCount", currentFloorCount)
	c.String(http.StatusOK,
		fmt.Sprintf("Success Create Floor [%d] Server[%d] Port[%d]", currentFloorCount, serverId, port))
}

// floor를 삭제하는 로직 ( Refresh 또는 맵 구조 변경용?
func deleteFloor(c *gin.Context) {
	ServerIDString := c.Param("serverId")
	serverId, err := strconv.Atoi(ServerIDString)
	if err != nil {
		c.String(http.StatusBadRequest,
			fmt.Sprintf("Invalid Parameter Delete Floor ServerId %d", serverId))
		return
	}

	delete(serverMap, serverId)
	c.String(http.StatusOK, fmt.Sprintf("Delete Floor Floor %d",  serverId))
	currentFloorCount -= 1
}