package routers

import (
	"dockerapi/app/api/v1"
	"dockerapi/global"
)

func init() {

	v := new(v1.UserRouter)
	api := global.GvaGinEngine.Group("auth")
	{

		api.POST("login", v.Login)
		api.POST("register", v.Register)

	}

}
