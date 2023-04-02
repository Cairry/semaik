package sdk

import (
	"context"
	"github.com/docker/docker/client"
	"log"
)

var (
	DockerCtx    context.Context
	DockerClient *client.Client
	err          error
)

func init() {

	DockerCtx = context.Background()

	// 远程连接
	DockerClient, err = client.NewClientWithOpts(client.WithHost("tcp://192.168.1.158:2376"),
		client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}

	// 本地连接
	//cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	//if err != nil {
	//	log.Fatal(err)
	//}

}
