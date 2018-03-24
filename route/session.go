package route

import "github.com/gin-gonic/gin"
import "net/http"

// 유저 세션을 만들어주는 역할
func getUserAccessInfo(c *gin.Context) {
	user := c.Param("user")

	// Redis에서 그냥 Sorting해서 해야하나? 미리 넣어놔야하나?
	// 게임 서버 내에서 게임 서버 접속 정보를 라우팅 하는 기능이 필요

	// Server -> User Connect Info
	// Logging으로 바꿀꺼
	c.String(http.StatusOK, "Hello %s", user)
}