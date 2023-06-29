package routers

import (
	docker2 "dockerapi/app/api/v1/docker"
	"dockerapi/global"
)

func init() {

	//.Use(middleware.JwtAuth(service.AppGuardName))
	api := global.GvaGinEngine.Group("api/v1")
	{

		// 容器
		api.GET("containers/list", docker2.ContainersList)
		api.GET("containers/state/:id", docker2.ContainerState)
		api.GET("containers/exec", docker2.ContainerWsSsh)
		api.POST("containers/create", docker2.ContainerCreate)
		api.POST("containers/option", docker2.ContainerOptions)
		api.POST("containers/log", docker2.ContainerLogs)

		// 镜像
		api.GET("images/list", docker2.ImagesList)
		api.POST("images/delete", docker2.ImageDelete)
		api.POST("images/pull", docker2.ImagePull)

		// 网络
		api.GET("networks/list", docker2.NetworksList)
		api.POST("networks/create", docker2.NetworksCreate)
		api.POST("networks/delete", docker2.NetworkDelete)

		// 卷
		api.GET("volumes/list", docker2.VolumesList)
		api.POST("volumes/create", docker2.VolumesCreate)
		api.POST("volumes/delete", docker2.VolumeDelete)

	}

}
