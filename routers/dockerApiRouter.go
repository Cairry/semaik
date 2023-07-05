package routers

import (
	"dockerapi/app/api/v1/docker"
	_ "dockerapi/docs"
	"dockerapi/global"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	d docker.DockerNodeApi
	c docker.ContainersApi
	i docker.ImagesApi
	n docker.NetworksApi
	v docker.VolumesApi
)

func init() {

	//.Use(middleware.JwtAuth(sdk.AppGuardName))
	api := global.GvaGinEngine.Group("api/v1")
	{

		// swagger
		api.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// docker
		node(api)
		cApi := api.Group("clouds/node/:nodeId")
		{
			container(cApi)
		}

	}

}

func node(ginEngine *gin.RouterGroup) {
	ginEngine.POST("clouds/node/create", d.Create)
	ginEngine.PUT("clouds/node/update", d.Update)
	ginEngine.DELETE("clouds/node/delete", d.Delete)
}

func container(ginEngine *gin.RouterGroup) {
	// 容器
	ginEngine.GET("containers/list", c.List)
	ginEngine.POST("containers/search", c.Search)
	ginEngine.GET("containers/state/:id", c.State)
	ginEngine.GET("containers/exec", c.Ssh)
	ginEngine.POST("containers/create", c.Create)
	ginEngine.POST("containers/option", c.Options)
	ginEngine.POST("containers/log", c.Logs)

	// 镜像
	ginEngine.GET("images/list", i.List)
	ginEngine.POST("images/search", i.Search)
	ginEngine.DELETE("images/delete", i.Delete)
	ginEngine.POST("images/pull", i.Pull)

	// 网络
	ginEngine.GET("networks/list", n.List)
	ginEngine.POST("networks/search", n.Search)
	ginEngine.POST("networks/create", n.Create)
	ginEngine.DELETE("networks/delete", n.Delete)

	// 卷
	ginEngine.GET("volumes/list", v.List)
	ginEngine.POST("volumes/search", v.Search)
	ginEngine.POST("volumes/create", v.Create)
	ginEngine.DELETE("volumes/delete", v.Delete)
}
