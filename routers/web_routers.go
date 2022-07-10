package routers

import (
	"gowebsocket/controllers/user"

	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.Engine) {
	// 用户组
	userRouter := router.Group("user")
	{
		userRouter.GET("/list", user.List)
		userRouter.GET("/online", user.Online)
		// userRouter.POST("/sendMessage", user.SendMessage)
		userRouter.POST("/sendMessageAll", user.SendMessageAll)
	}
}
