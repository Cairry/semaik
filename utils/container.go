package utils

import (
	"context"
	"dockerapi/global"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"strings"
)

// CalculateCPUPercentUnix 计算cpu
func CalculateCPUPercentUnix(previousCPU, previousSystem uint64, v *types.StatsJSON) float64 {
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

// CalculateBlockIO 计算io
func CalculateBlockIO(blkio types.BlkioStats) (blkRead float64, blkWrite float64) {
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

// CalculateNetwork 计算网卡
func CalculateNetwork(network map[string]types.NetworkStats) (float64, float64) {
	var rx, tx float64

	for _, v := range network {
		rx += float64(v.RxBytes) / 1024
		tx += float64(v.TxBytes) / 1024
	}
	return rx, tx
}

// ExecContainer 连接 Container
func ExecContainer(ctx context.Context, containerID string, user string, command []string) (hr types.HijackedResponse, err error) {

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
	exec, err := global.GvaDockerCli.ContainerExecCreate(ctx, containerID, execConfig)
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
	resp, err := global.GvaDockerCli.ContainerExecAttach(ctx, exec.ID, attachConfig)
	if err != nil {
		log.Println("连接 Container 错误:", err)
		return
	}

	return resp, nil

}

// WsWriterCopy 读取终端返回的数据, 写入到 ws 中
func WsWriterCopy(reader io.Reader, writer *websocket.Conn) {

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

// WsReaderCopy 从 ws 读取用户输入的数据, 写入到终端中
func WsReaderCopy(reader *websocket.Conn, writer io.Writer) {

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

// FormatFileSize 计算镜像大小
func FormatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else {
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}
