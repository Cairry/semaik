package router

import (
	"dockerapi/global"
	"dockerapi/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func init() {

	global.GvaGinEngine.GET("ping", middleware.Auth(), func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"msg": "pong",
		})
	})

}
