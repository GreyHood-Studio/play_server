package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"github.com/GreyHood-Studio/play_server/controller"
)

func healthCheck(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf("Health Play Server"))
}

// router package는 gameserver들을 컨트롤 하기 위한 로직
func SetAPIRoute(router *gin.Engine, maxRoomCount int) {
	// 서버의 상태를 가지고 오는 정보
	//router.GET("/floors", getServerStatus)
	router.GET("/ping", healthCheck)

	// Floor 관리
	roomGroup := router.Group("/room")
	{
		roomCtrl := controller.NewRoomController(maxRoomCount)
		roomGroup.GET("/:serverID", roomCtrl.GetRoom)
		roomGroup.POST("/:serverID", roomCtrl.CreateRoom)
		roomGroup.DELETE("/:serverID", roomCtrl.DeleteRoom)
	}

	// User 강퇴 및 해당 User가 접속해 있는지 확인하기 위한 로직
	//userGroup := router.Group("/user")
	//{
	//
	//}
}