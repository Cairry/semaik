package router

import (
	"dockerapi/app/service"
	"dockerapi/global"
	"dockerapi/middleware"
	"dockerapi/routers/api/docker"
	"github.com/gin-gonic/gin"
)

func init() {

	api := global.GvaGinEngine.Group("api").Use(middleware.JwtAuth(service.AppGuardName))
	{
		api.GET("test", func(context *gin.Context) {
			context.JSON(200, gin.H{"msg": "success"})
		})

		// 容器
		api.GET("containers/list", docker.ContainersList)
		api.POST("containers/create", docker.ContainerCreate)
		api.POST("containers/option", docker.ContainerOptions)

		// 镜像
		api.GET("images/list", docker.ImagesList)
		api.POST("images/delete", docker.ImageDelete)
		api.POST("images/pull", docker.ImagePull)

		// 网络
		api.GET("networks/list", docker.NetworksList)
		api.POST("networks/create", docker.NetworksCreate)
		api.POST("networks/delete", docker.NetworkDelete)

		// 卷
		api.GET("volumes/list", docker.VolumesList)
		api.POST("volumes/create", docker.VolumesCreate)
		api.POST("volumes/delete", docker.VolumeDelete)
	}

}
