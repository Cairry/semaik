package router

import (
	"dockerapi/global"
	"dockerapi/routers/api/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func init() {

	api := global.GvaGinEngine.Group("auth")
	{
		api.GET("info", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{
				"msg": "success",
			})
		})

		api.POST("login", auth.Login)
		api.POST("register", auth.Register)

	}

}
