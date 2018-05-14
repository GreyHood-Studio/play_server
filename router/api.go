package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
)

func healthCheck(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf("Health Play Server"))
}

// router package는 gameserver들을 컨트롤 하기 위한 로직
func SetAPIRoute(router *gin.Engine) {
	// 서버의 상태를 가지고 오는 정보
	//router.GET("/floors", getServerStatus)
	router.GET("/ping", healthCheck)

	// Floor 관리
	floor := router.Group("/room")
	{
		floor.GET("/:serverID", getRoom)
		floor.POST("/:serverID", createRoom)
		floor.DELETE("/:serverID", deleteRoom)
	}

	// 유저 정보 입력 삭제, 세션 종료, 강제 퇴장등의 기능
	user := router.Group("/user")
	{
		// 유저가 접속할 ip와 port번호를 전달
		user.GET("/:userID", getUserAccessInfo)
	}
}