package controller

import (
	"github.com/GreyHood-Studio/play_server/model"
	"strconv"
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"github.com/GreyHood-Studio/play_server/network"
)

// controller 객체들은 모두 router로부터 설정된 함수들을 받아서 처리하는 함수들이다.
//UserController ...
type RoomController struct{
	// key == server id
	MaxRoomCount	int
	// client들이 게임 room object에 접근하기 위한 변수
	// roomMap Key값은 serverId (client들은 본인들의 서버 id를 알고있음
	RoomServers		map[int]*RoomServer
}

type RoomServer struct {
	NetworkServer	*network.GameServer
	Room			*model.Room
}

// Floor의 정보를 가져오는 로직
// 만약에 Redis에 있으면 필요 없음
func (ctrl RoomController)GetRoom(c *gin.Context) {
	ServerIDString := c.Param("serverID")
	serverId, _ := strconv.Atoi(ServerIDString)

	// map 키 체크
	_, exists := ctrl.RoomServers[serverId]
	if exists {

	} else {
		c.String(http.StatusNotFound, fmt.Sprintf("Don't Exist ServerId %d", serverId))
		return
	}
}

func (ctrl RoomController) CreateRoom(c *gin.Context) {
	// 실제 게임 서버의 층 오픈
	// Packet 층 말고도, Floor의 정보를 보관하는 로직이 필요
	ServerIDString := c.PostForm("serverId")
	portString := c.PostForm("port")
	maxConnString := c.PostForm("maxConn")
	mapIDString := c.PostForm("mapId")
	port, pErr := strconv.Atoi(portString)
	serverId, sErr := strconv.Atoi(ServerIDString)
	maxConn, mErr := strconv.Atoi(maxConnString)
	mapId, mErr := strconv.Atoi(mapIDString)

	if pErr != nil || sErr != nil || mErr != nil {
		c.String(http.StatusBadRequest,
			fmt.Sprintf("Invalid Parameter create Floor Server[%d] Port[%d]", serverId, port))
		return
	}
	fmt.Println(port, serverId, maxConn)

	// server create validation logic
	if len(ctrl.RoomServers) >= ctrl.MaxRoomCount {
		c.String(http.StatusNotAcceptable, fmt.Sprintf("Already Max Server Count"))
		return
	}

	_, exists := ctrl.RoomServers[serverId]
	//if 포트 또는 서버아이디가 중복된 경우 Return
	if exists {
		c.String(http.StatusNotAcceptable,
			fmt.Sprintf("Already Exist Server Port %d OR ServerId %d", port, serverId))
		return
	}
	// server validation logic complete

	// create room and server
	room := model.NewRoom(maxConn, mapId)
	netServer := network.NewServer(serverId, port, maxConn, room)

	go netServer.Run()

	roomServer := &RoomServer{NetworkServer:netServer, Room:room}
	// 서버 map에 serverID를 기준으로 서버를 생성
	ctrl.RoomServers[serverId] = roomServer

	// Redis에 현재 생성한 룸을 저장

	// 현재 생성된 서버의 수
	fmt.Printf("currentFloorCount %d\n", len(ctrl.RoomServers))
	c.String(http.StatusOK,
		fmt.Sprintf("Success Create Room MapId[%d] Server[%d] Port[%d]", mapId, serverId, port))
}

func (ctrl RoomController) DeleteRoom(c *gin.Context) {
	ServerIDString := c.Param("serverID")
	serverId, err := strconv.Atoi(ServerIDString)
	if err != nil {
		c.String(http.StatusBadRequest,
			fmt.Sprintf("Invalid Parameter Delete Room ServerId %d", serverId))
		return
	}
	delete(ctrl.RoomServers, serverId)

	// Redis에 현재 삭제한 룸을 삭제

	c.String(http.StatusOK, fmt.Sprintf("Delete Room Floor %d",  serverId))
}

func NewRoomController(maxRoomCount int) *RoomController{
	controller := &RoomController{
		RoomServers: make(map[int]*RoomServer),
		MaxRoomCount: maxRoomCount,
	}

	return controller
}