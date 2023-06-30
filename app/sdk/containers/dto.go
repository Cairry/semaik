package containers

import (
	"time"
)

// ContainerListStruct 容器列表
type ContainerListStruct struct {
	ID      string
	Name    string
	Image   string
	ImageID string
	Ports   []string
	Created time.Time
	State   string
	Status  string
}

// ContainerCreateStruct 创建容器
type ContainerCreateStruct struct {
	Name               string            `json:"name"`
	Image              string            `json:"image"`
	Env                []string          `json:"env"`
	Cmd                []string          `json:"cmd"`
	Labels             map[string]string `json:"labels"`
	EnableResourceList bool              `json:"enableResourceList"`
	NanoCPUs           int64             `json:"nanoCPUs"`
	Memory             int64             `json:"memory"`
	EnablePublicPort   bool              `json:"enablePublicPort"`
	ExposedPorts       []struct {
		ContainerPort int `json:"containerPort"`
		HostPort      int `json:"hostPort"`
	} `json:"exposedPorts"`
	EnableVolumeMount bool `json:"enableVolumeMount"`
	Volumes           []struct {
		SourceDir    string `json:"sourceDir"`
		ContainerDir string `json:"containerDir"`
		Mode         string `json:"mode"`
	} `json:"volumes"`
}

type ContainerOperationStruct struct {
	Name      string `json:"name" validate:"required"`
	Operation string `json:"operation" validate:"required,oneof=start stop restart kill pause unpause rename remove"`
	NewName   string `json:"newName"`
}

type ContainerLogsStruct struct {
	ContainerName string `json:"containerName"`
	ContainerID   string `json:"containerID"`
	IsWatch       bool   `json:"isWatch"`
	Mode          string `json:"mode"`
}

type ContainerStats struct {
	CPUPercent float64   `json:"cpuPercent"`
	Memory     float64   `json:"memory"`
	Cache      float64   `json:"cache"`
	IORead     float64   `json:"ioRead"`
	IOWrite    float64   `json:"ioWrite"`
	NetworkRX  float64   `json:"networkRX"`
	NetworkTX  float64   `json:"networkTX"`
	ShotTime   time.Time `json:"shotTime"`
}

type ContainerWsSshStruct struct {
	ContainerID string   `json:"containerID"`
	User        string   `json:"user"`
	Command     []string `json:"command"`
}
