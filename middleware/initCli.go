package middleware

import (
	"dockerapi/app/service/docker/node"
	"dockerapi/global"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func InitCli(ctx *gin.Context) {

	id := ctx.Param("nodeId")

	var nodeInfo node.DockerNode
	global.GvaDatabase.Where("name = ?", id).First(&nodeInfo)

	// 远程连接
	cli, err := client.NewClientWithOpts(client.WithHost(nodeInfo.Host), client.WithAPIVersionNegotiation(), client.WithTimeout(3*time.Second))
	if err != nil {
		log.Fatal("连接 Docker 失败:", err)
	}

	// 本地连接
	//cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	//if err != nil {
	//	log.Fatal(err)
	//}

	global.GvaDockerCli = cli

}
