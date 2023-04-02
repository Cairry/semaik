package middleware

import (
	"dockerapi/app/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Auth() gin.HandlerFunc {

	return func(context *gin.Context) {
		token := context.Request.Header.Get("token")
		if token != "zux" {
			response.Response(context, http.StatusOK, 200, nil, "用户未登陆")
		}
		context.Abort()
	}
}
