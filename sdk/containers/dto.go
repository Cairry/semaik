package containers

import "time"

// ContainerListStruct 容器列表
type ContainerListStruct struct {
	ID      string
	Names   string
	Image   string
	ImageID string
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