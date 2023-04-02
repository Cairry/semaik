package service

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type jwtService struct {
}

var JwtService = new(jwtService)

// JwtUser 获取用户信息
type JwtUser interface {
	GetUid() string
}

type CustomClaims struct {
	jwt.StandardClaims
}

const (
	// TokenType Token 类型
	TokenType = "bearer"
	// AppGuardName 颁发者
	AppGuardName = "app"
	// Expire 失效时间
	Expire = 40000
)

type TokenOutPut struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
	TokenType   string `json:"tokenType"`
}

// CreateToken 创建 Token
func (jwtService *jwtService) CreateToken(GuardName string, user JwtUser) (tokenData TokenOutPut, err error, token *jwt.Token) {

	token = jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		CustomClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Unix() + Expire,
				Id:        user.GetUid(),
				Issuer:    GuardName,
			},
		},
	)

	tokenStr, err := token.SignedString([]byte("zux"))
	tokenData = TokenOutPut{
		tokenStr,
		int(Expire),
		TokenType,
	}
	return

}
