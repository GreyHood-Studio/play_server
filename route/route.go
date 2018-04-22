package route

import "github.com/gin-gonic/gin"


// route package는 gameserver들을 컨트롤 하기 위한 로직
func SetAPIRoute(router *gin.Engine) {
	// 서버의 상태를 가지고 오는 정보
	//router.GET("/floors", getServerStatus)

	// Simple group: v1
	floor := router.Group("/floor")
	{
		floor.GET("/:serverID", getFloor)
		floor.POST("/:serverID", createFloor)
		floor.DELETE("/:serverID", deleteFloor)
	}

	// 유저 정보 입력 삭제, 세션 종료, 강제 퇴장등의 기능
	user := router.Group("/user")
	{
		// 유저가 접속할 ip와 port번호를 전달
		user.GET("/:userID", getUserAccessInfo)
	}

	// 몬스터, 게임 내에 설정 정보 변경등에 대한 내용
	gameManager := router.Group("/manager")
	{
		gameManager.GET("/:serverID", regenMonster)
	}
}