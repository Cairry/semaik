package node

import (
	"dockerapi/global"
	"time"
)

type DockerNodeService struct{}

/*
	Create
	创建 Docker 节点
*/
func (dns DockerNodeService) Create(info DockerNode) (interface{}, error) {

	info.Created = time.Now()
	err := global.GvaDatabase.Create(&info).Error
	if err != nil {
		return nil, err
	}

	return info, nil
}

/*
	Update
	更新 Docker 节点信息
*/
func (dns DockerNodeService) Update(info DockerNode) (interface{}, error) {

	var newInfo DockerNode
	err := global.GvaDatabase.Model(info).Where("name = ?", info.Name).Update("host", info.Host).First(&newInfo).Error
	if err != nil {
		return nil, err
	}
	return newInfo, nil

}

/*
	Delete
	删除 Docker 节点信息
*/
func (dns DockerNodeService) Delete(info DockerNode) (interface{}, error) {

	var newInfo []DockerNode
	err := global.GvaDatabase.Where("name = ?", info.Name).Find(&newInfo).Delete(info).Error
	if err != nil {
		return nil, err
	}
	return newInfo, nil

}
