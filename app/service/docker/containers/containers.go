package containers

import (
	"context"
	"dockerapi/app/service/docker/images"
	"dockerapi/global"
	"dockerapi/utils"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

type ContainerService struct{}

/*
	List
	获取容器列表
*/
func (cs ContainerService) List(ctx context.Context) []ContainerListStruct {

	// 获取容器列表

	containers, err := global.GvaDockerCli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		log.Fatal(err)
	}

	// 初始化一个容器列表切片，存放容器列表数据
	result := make([]ContainerListStruct, 0, len(containers))

	for _, container := range containers {

		var ports []string
		for _, port := range container.Ports {
			if port.IP == "::" || port.PublicPort == 0 {
				continue
			}
			tmp := fmt.Sprintf("%d:%d/%s", port.PublicPort, port.PrivatePort, port.Type)
			ports = append(ports, tmp)
		}

		// 读取容器列表数据
		tmp := ContainerListStruct{
			ID:      container.ID[:10],
			Name:    container.Names[0][1:],
			Image:   container.Image,
			ImageID: strings.Split(container.ImageID, ":")[1][:10],
			Ports:   ports,
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
	Search
	获取指定容器
*/
func (cs ContainerService) Search(ctx context.Context, info string) []ContainerListStruct {

	containers, err := global.GvaDockerCli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	var containerList []types.Container

	if len(info) != 0 {
		length, count := len(containers), 0
		for count < length {
			for _, v := range containers[count].Names {
				if strings.Contains(v, info) {
					containerList = append(containerList, containers[count])
					break
				}
			}
			count++
		}
	} else {
		containerList = containers
	}

	var result []ContainerListStruct

	for _, container := range containerList {
		var ports []string
		for _, port := range container.Ports {
			if port.IP == "::" || port.PublicPort == 0 {
				continue
			}
			tmp := fmt.Sprintf("%d:%d/%s", port.PublicPort, port.PrivatePort, port.Type)
			ports = append(ports, tmp)
		}

		// 读取容器列表数据
		tmp := ContainerListStruct{
			ID:      container.ID[:10],
			Name:    container.Names[0][1:],
			Image:   container.Image,
			ImageID: strings.Split(container.ImageID, ":")[1][:10],
			Ports:   ports,
			Created: time.Unix(container.Created, 0),
			State:   container.State,
			Status:  container.Status,
		}

		// 将容器列表数据写入到初始化的容器列表的切片中
		result = append(result, tmp)
	}

	return result
}

/*
	Create
	创建容器
*/
func (cs ContainerService) Create(ctx context.Context, req ContainerCreateStruct) error {

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
		hostConfig container.HostConfig
	)

	// 判断是否开启资源限制
	if req.EnableResourceList {
		hostConfig.NanoCPUs = req.NanoCPUs * 1000000000
		hostConfig.Memory = req.Memory
	}

	// 判断是否暴露端口
	if req.EnablePublicPort {
		hostConfig.PortBindings = make(nat.PortMap)
		for _, port := range req.ExposedPorts {
			bindItem := nat.PortBinding{HostPort: strconv.Itoa(port.HostPort)}
			hostConfig.PortBindings[nat.Port(fmt.Sprintf("%d/tcp", port.ContainerPort))] = []nat.PortBinding{bindItem}
		}
	}

	// 判断是否配置挂载点
	if req.EnableVolumeMount {
		config.Volumes = make(map[string]struct{})
		for _, volume := range req.Volumes {
			config.Volumes[volume.ContainerDir] = struct{}{}
			hostConfig.Binds = append(hostConfig.Binds, fmt.Sprintf("%s:%s:%s", volume.SourceDir, volume.ContainerDir, volume.Mode))
		}
	}

	// 拉取镜像
	err := images.ImageService{}.Pull(ctx, req.Image)
	if err != nil {
		log.Println("Image pull failed:", err)
		return err
	}

	// 创建容器
	resp, err := global.GvaDockerCli.ContainerCreate(ctx, config, &hostConfig, nil, nil, req.Name)
	if err != nil {
		log.Println("Container Create failed:", err)
		return err
	}

	// 启动容器
	err = global.GvaDockerCli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Println("Container Start failed:", err)
		return err
	}

	return nil
}

/*
	Options
	容器操作选项
*/
func (cs ContainerService) Options(ctx context.Context, req ContainerOperationStruct) error {

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
		err = global.GvaDockerCli.ContainerStart(ctx, req.Name, types.ContainerStartOptions{})
	case "stop":
		err = global.GvaDockerCli.ContainerStop(ctx, req.Name, nil)
	case "restart":
		err = global.GvaDockerCli.ContainerRestart(ctx, req.Name, nil)
	case "kill":
		err = global.GvaDockerCli.ContainerKill(ctx, req.Name, "SIGKILL")
	case "pause":
		err = global.GvaDockerCli.ContainerPause(ctx, req.Name)
	case "unpause":
		err = global.GvaDockerCli.ContainerUnpause(ctx, req.Name)
	case "rename":
		err = global.GvaDockerCli.ContainerRename(ctx, req.Name, req.NewName)
	case "remove":
		err = global.GvaDockerCli.ContainerRemove(ctx, req.Name, options)

	}
	return err
}

/*
	Logs
	获取容器日志
*/
func (cs ContainerService) Logs(ctx context.Context, req ContainerLogsStruct) (string, error) {

	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false,
	}

	if req.Mode != "all" {
		options.Since = req.Mode
	}

	logs, err := global.GvaDockerCli.ContainerLogs(ctx, req.ContainerName, options)
	if err != nil {
		log.Println(err)
	}

	body, err := ioutil.ReadAll(logs)
	if err != nil {
		log.Println(err)
	}
	return string(body), nil

}

/*
	State
	获取容器资源利用率
*/
func (cs ContainerService) State(ctx context.Context, id string) (*ContainerStats, error) {

	res, err := global.GvaDockerCli.ContainerStats(ctx, id, false)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var (
		stats *types.StatsJSON
		data  ContainerStats
	)
	if err := json.Unmarshal(body, &stats); err != nil {
		return nil, err
	}

	preCpu := stats.PreCPUStats.CPUUsage.TotalUsage
	preSystem := stats.PreCPUStats.SystemUsage
	data.CPUPercent = utils.CalculateCPUPercentUnix(preCpu, preSystem, stats)
	data.IORead, data.IOWrite = utils.CalculateBlockIO(stats.BlkioStats)
	data.Memory = float64(stats.MemoryStats.Usage) / 1024 / 1024
	if cache, ok := stats.MemoryStats.Stats["cache"]; ok {
		data.Cache = float64(cache) / 1024 / 1024
	}
	data.Memory = data.Memory - data.Cache
	data.NetworkRX, data.NetworkTX = utils.CalculateNetwork(stats.Networks)
	data.ShotTime = stats.Read
	return &data, nil

}

/*
	Ssh
	连接容器终端
*/
func (cs ContainerService) Ssh(ctx context.Context, req ContainerWsSshStruct, c *gin.Context, wsConn *websocket.Conn) (err error) {

	// 连接容器
	execResp, err := utils.ExecContainer(ctx, req.ContainerID, req.User, req.Command)
	if err != nil {
		log.Println("请求失败")
		return err
	}

	defer execResp.Close()

	defer func() {
		execResp.Conn.Write([]byte("exit\r"))
	}()

	// 开启协程实时监听终端反馈的数据
	go func() {
		utils.WsWriterCopy(execResp.Conn, wsConn)
	}()

	// 程序阻塞，读取 ws 数据写入终端
	utils.WsReaderCopy(wsConn, execResp.Conn)

	return nil

}
