package routers

import (
	"dockerapi/app/api/v1/user"
	"dockerapi/global"
)

func init() {

	var (
		v user.UserRouter
	)
	api := global.GvaGinEngine.Group("auth")
	{

		api.POST("login", v.Login)
		api.POST("register", v.Register)

	}

}
