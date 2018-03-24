package route

import "github.com/gin-gonic/gin"

func SetAPIRoute(router *gin.Engine) {

	// 유저가 접속할 ip와 port번호를 전달
	router.GET("/join/:user", getUserAccessInfo)
	// 서버의 상태를 가지고 오는 정보
	router.GET("/status", getServerStatus)

	// Simple group: v1
	floor := router.Group("/floor")
	{
		floor.GET("/status", getFloorStatus)
		floor.POST("/createFloor", createFloor)
		floor.DELETE("/deleteFloor/:floorid", createFloor)
	}
}