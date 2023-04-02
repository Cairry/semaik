package middleware

import (
	"dockerapi/app/response"
	"dockerapi/app/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JwtAuth(GuardNmae string) gin.HandlerFunc {

	return func(context *gin.Context) {
		// 获取 Token
		tokenStr := context.Request.Header.Get("Authorization")
		if tokenStr == "" {
			response.TokenFail(context)
			context.Abort()
			return
		}
		// Bearer Token, 获取 Token 值
		tokenStr = tokenStr[len(service.TokenType)+1:]

		// 校验 Token
		token, err := jwt.ParseWithClaims(tokenStr, &service.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("zux"), nil
		})
		if err != nil {
			response.TokenFail(context)
			context.Abort()
			return
		}

		// 断言
		claims := token.Claims.(*service.CustomClaims)

		// 发布者校验
		if claims.Issuer != GuardNmae {
			response.TokenFail(context)
			context.Abort()
			return
		}

		context.Set("token", token)
		context.Set("id", claims.Id)

	}

}
