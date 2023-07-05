package node

import "time"

type DockerNode struct {
	Name    string    `json:"name" gorm:"unique"`
	Host    string    `json:"host"`
	Created time.Time `json:"created"`
}
