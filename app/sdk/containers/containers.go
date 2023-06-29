package containers

import (
	"dockerapi/app/sdk"
	"dockerapi/app/sdk/images"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
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
	resp, err := sdk.DockerClient.ContainerCreate(sdk.DockerCtx, config, &hostConfg, nil, nil, req.Name)
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

/*
	ContainerLogs
	获取容器日志
*/
func ContainerLogs(req ContainerLogsStruct) (string, error) {

	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false,
	}

	if req.Mode != "all" {
		options.Since = req.Mode
	}

	logs, err := sdk.DockerClient.ContainerLogs(sdk.DockerCtx, req.ContainerName, options)
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
	ContainerState
	获取容器资源利用率
*/
func ContainerState(id string) (*ContainerStats, error) {

	res, err := sdk.DockerClient.ContainerStats(sdk.DockerCtx, id, false)
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
	data.CPUPercent = calculateCPUPercentUnix(preCpu, preSystem, stats)
	data.IORead, data.IOWrite = calculateBlockIO(stats.BlkioStats)
	data.Memory = float64(stats.MemoryStats.Usage) / 1024 / 1024
	if cache, ok := stats.MemoryStats.Stats["cache"]; ok {
		data.Cache = float64(cache) / 1024 / 1024
	}
	data.Memory = data.Memory - data.Cache
	data.NetworkRX, data.NetworkTX = calculateNetwork(stats.Networks)
	data.ShotTime = stats.Read
	return &data, nil

}

// 计算cpu
func calculateCPUPercentUnix(previousCPU, previousSystem uint64, v *types.StatsJSON) float64 {
	var (
		cpuPercent  = 0.0
		cpuDelta    = float64(v.CPUStats.CPUUsage.TotalUsage) - float64(previousCPU)
		systemDelta = float64(v.CPUStats.SystemUsage) - float64(previousSystem)
	)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(v.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}
	return cpuPercent
}

// 计算io
func calculateBlockIO(blkio types.BlkioStats) (blkRead float64, blkWrite float64) {
	for _, bioEntry := range blkio.IoServiceBytesRecursive {
		switch strings.ToLower(bioEntry.Op) {
		case "read":
			blkRead = (blkRead + float64(bioEntry.Value)) / 1024 / 1024
		case "write":
			blkWrite = (blkWrite + float64(bioEntry.Value)) / 1024 / 1024
		}
	}
	return
}

// 计算网卡
func calculateNetwork(network map[string]types.NetworkStats) (float64, float64) {
	var rx, tx float64

	for _, v := range network {
		rx += float64(v.RxBytes) / 1024
		tx += float64(v.TxBytes) / 1024
	}
	return rx, tx
}

/*
	ContainerSsh
	连接容器终端
*/

func execContainer(containerID string, user string, command []string) (hr types.HijackedResponse, err error) {

	//exec 配置
	execConfig := types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		User:         user,
		Cmd:          command,
	}

	//创建 exec 实例
	exec, err := sdk.DockerClient.ContainerExecCreate(sdk.DockerCtx, containerID, execConfig)
	if err != nil {
		log.Println("创建 exec 实例错误:", err)
		return
	}

	// Attach 配置
	attachConfig := types.ExecStartCheck{
		Tty:    true,
		Detach: false,
	}

	// 连接 Container
	resp, err := sdk.DockerClient.ContainerExecAttach(sdk.DockerCtx, exec.ID, attachConfig)
	if err != nil {
		log.Println("连接 Container 错误:", err)
		return
	}

	return resp, nil

}

func ContainerSsh(req ContainerWsSshStruct, c *gin.Context, wsConn *websocket.Conn) (err error) {

	// 连接容器
	execResp, err := execContainer(req.ContainerID, req.User, req.Command)
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
		wsWriterCopy(execResp.Conn, wsConn)
	}()

	// 程序阻塞，读取 ws 数据写入终端
	wsReaderCopy(wsConn, execResp.Conn)

	return nil

}

// 读取终端返回的数据, 写入到 ws 中
func wsWriterCopy(reader io.Reader, writer *websocket.Conn) {

	buf := make([]byte, 8192)
	for {
		nr, err := reader.Read(buf)
		//fmt.Printf("终端输出的数据 ---> %s", buf[:nr])
		if nr > 0 {
			err := writer.WriteMessage(websocket.BinaryMessage, buf[:nr])
			if err != nil {
				return
			}
		}
		if err != nil {
			return
		}
	}

}

// 从 ws 读取用户输入的数据, 写入到终端中
func wsReaderCopy(reader *websocket.Conn, writer io.Writer) {

	for {
		messageType, p, err := reader.ReadMessage()
		if err != nil {
			return
		}

		//fmt.Printf("用户输入的数据 ---> %s\n", p)
		if messageType == websocket.TextMessage {
			writer.Write(p)
		}
	}

}
