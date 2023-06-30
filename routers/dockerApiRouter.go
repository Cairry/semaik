package routers

import (
	docker2 "dockerapi/app/api/v1/docker"
	_ "dockerapi/docs"
	"dockerapi/global"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {

	var (
		c docker2.ContainersApi
		i docker2.ImagesApi
		n docker2.NetworksApi
		v docker2.VolumesApi
	)

	//.Use(middleware.JwtAuth(service.AppGuardName))
	api := global.GvaGinEngine.Group("api/v1")
	{

		// swagger
		api.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// 容器
		api.GET("containers/list", c.List)
		api.POST("containers/search", c.Search)
		api.GET("containers/state/:id", c.State)
		api.GET("containers/exec", c.Ssh)
		api.POST("containers/create", c.Create)
		api.POST("containers/option", c.Options)
		api.POST("containers/log", c.Logs)

		// 镜像
		api.GET("images/list", i.List)
		api.POST("images/search", i.Search)
		api.DELETE("images/delete", i.Delete)
		api.POST("images/pull", i.Pull)

		// 网络
		api.GET("networks/list", n.List)
		api.POST("networks/search", n.Search)
		api.POST("networks/create", n.Create)
		api.DELETE("networks/delete", n.Delete)

		// 卷
		api.GET("volumes/list", v.List)
		api.POST("volumes/search", v.Search)
		api.POST("volumes/create", v.Create)
		api.DELETE("volumes/delete", v.Delete)

	}

}
