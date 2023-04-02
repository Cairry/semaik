package containers

import (
	"dockerapi/sdk"
	"dockerapi/sdk/images"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"log"
	"strconv"
	"strings"
	"time"
)

/*
	ContainersList
	获取容器列表
*/
func ContainersList() []ContainerListStruct {

	// 获取容器列表
	containers, err := sdk.DockerClient.ContainerList(sdk.DockerCtx, types.ContainerListOptions{All: true})
	if err != nil {
		log.Fatal(err)
	}

	// 初始化一个容器列表切片，存放容器列表数据
	result := make([]ContainerListStruct, 0, len(containers))

	for _, container := range containers {

		// 读取容器列表数据
		tmp := ContainerListStruct{
			ID:      container.ID[:10],
			Names:   container.Names[0][1:],
			Image:   container.Image,
			ImageID: strings.Split(container.ImageID, ":")[1][:10],
			Created: time.Unix(container.Created, 0),
			State:   container.State,
			Status:  container.Status,
		}

		// 将容器列表数据写入到初始化的容器列表的切片中
		result = append(result, tmp)

	}

	// 打印容器列表
	return result
}

/*
	ContainerCreate
	创建容器
*/
func ContainerCreate(req ContainerCreateStruct) error {

	// 定义创建容器时的配置信息
	config := &container.Config{
		Tty:       true,
		OpenStdin: true,
		Env:       req.Env,
		Cmd:       req.Cmd,
		Image:     req.Image,
		Labels:    req.Labels,
	}

	var (
		hostConfg container.HostConfig
	)

	// 判断是否开启资源限制
	if req.EnableResourceList {
		hostConfg.NanoCPUs = req.NanoCPUs * 1000000000
		hostConfg.Memory = req.Memory
	}

	// 判断是否暴露端口
	if req.EnablePublicPort {
		hostConfg.PortBindings = make(nat.PortMap)
		for _, port := range req.ExposedPorts {
			bindItem := nat.PortBinding{HostPort: strconv.Itoa(port.HostPort)}
			hostConfg.PortBindings[nat.Port(fmt.Sprintf("%d/tcp", port.ContainerPort))] = []nat.PortBinding{bindItem}
		}
	}

	// 判断是否配置挂载点
	if req.EnableVolumeMount {
		config.Volumes = make(map[string]struct{})
		for _, volume := range req.Volumes {
			config.Volumes[volume.ContainerDir] = struct{}{}
			hostConfg.Binds = append(hostConfg.Binds, fmt.Sprintf("%s:%s:%s", volume.SourceDir, volume.ContainerDir, volume.Mode))
		}
	}

	// 拉取镜像
	err := images.PullImage(req.Image)
	if err != nil {
		log.Println("Image pull failed:", err)
		return err
	}

	// 创建容器
	resp, err := sdk.DockerClient.ContainerCreate(sdk.DockerCtx, config, &hostConfg, nil, req.Name)
	if err != nil {
		log.Println("Container Create failed:", err)
		return err
	}

	// 启动容器
	err = sdk.DockerClient.ContainerStart(sdk.DockerCtx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Println("Container Start failed:", err)
		return err
	}

	return nil
}

/*
	ContainerOptions
	容器操作选项
*/
func ContainerOptions(req ContainerOperationStruct) error {

	var (
		err error
	)

	// 定义容器删除时的选项配置
	options := types.ContainerRemoveOptions{
		RemoveVolumes: false,
		RemoveLinks:   false,
		Force:         true,
	}

	// 获取操作选项并执行
	switch req.Operation {
	case "start":
		err = sdk.DockerClient.ContainerStart(sdk.DockerCtx, req.Name, types.ContainerStartOptions{})
	case "stop":
		err = sdk.DockerClient.ContainerStop(sdk.DockerCtx, req.Name, nil)
	case "restart":
		err = sdk.DockerClient.ContainerRestart(sdk.DockerCtx, req.Name, nil)
	case "kill":
		err = sdk.DockerClient.ContainerKill(sdk.DockerCtx, req.Name, "SIGKILL")
	case "pause":
		err = sdk.DockerClient.ContainerPause(sdk.DockerCtx, req.Name)
	case "unpause":
		err = sdk.DockerClient.ContainerUnpause(sdk.DockerCtx, req.Name)
	case "rename":
		err = sdk.DockerClient.ContainerRename(sdk.DockerCtx, req.Name, req.NewName)
	case "remove":
		err = sdk.DockerClient.ContainerRemove(sdk.DockerCtx, req.Name, options)

	}
	return err
}
