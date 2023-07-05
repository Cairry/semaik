package networks

import (
	"context"
	"dockerapi/global"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"log"
	"strings"
)

type NetworkService struct{}

/*
	List
	网络列表
*/
func (ns NetworkService) List(ctx context.Context) []NetworkList {

	networks, err := global.GvaDockerCli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// 初始化网络列表切片
	result := make([]NetworkList, 0, len(networks))

	// 循环读取
	for _, network := range networks {
		// 获取网络信息
		tmp := NetworkList{
			Name:    network.Name,
			ID:      network.ID[:10],
			Driver:  network.Driver,
			Scope:   network.Scope,
			Created: network.Created,
		}

		// 将网络信息加入到网络列表切片中
		result = append(result, tmp)

	}

	// 打印网络列表
	return result
}

/*
	Search
	搜索网络
*/
func (ns NetworkService) Search(ctx context.Context, info string) []NetworkList {

	networks, err := global.GvaDockerCli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	var networksList []types.NetworkResource

	if len(info) != 0 {
		length, count := len(networks), 0
		for count < length {
			if strings.Contains(networks[count].Name, info) {
				networksList = append(networksList, networks[count])
			}
			count++
		}
	} else {
		networksList = networks
	}

	result := make([]NetworkList, 0, len(networks))

	for _, network := range networksList {
		tmp := NetworkList{
			Name:    network.Name,
			ID:      network.ID[:10],
			Driver:  network.Driver,
			Scope:   network.Scope,
			Created: network.Created,
		}

		result = append(result, tmp)

	}

	return result

}

/*
	Create
	创建网络
*/
func (ns NetworkService) Create(ctx context.Context, req NetworkCreateStruct) error {

	// 先判断是否存在
	for _, v := range ns.List(ctx) {
		if v.Name == req.Name {
			return errors.New("无法创建网络, 它已存在")
		}
	}

	// 定义IPAM网络配置
	reqIPAMConfig := network.IPAMConfig{
		Subnet:     req.IPAM.Config[0].Subnet,
		IPRange:    req.IPAM.Config[0].IPRange,
		Gateway:    req.IPAM.Config[0].Gateway,
		AuxAddress: req.IPAM.Config[0].AuxAddress,
	}
	reqIPAM := network.IPAM{
		Driver:  req.IPAM.Driver,
		Options: req.IPAM.Options,
		Config:  []network.IPAMConfig{reqIPAMConfig},
	}

	// 定义基础创建配置
	options := types.NetworkCreate{
		Driver:  req.Driver,
		Options: req.Options,
		Labels:  req.Labels,
	}

	// 判断是否启用IPAM配置
	if req.EnableIPAM == true {
		options.IPAM = &reqIPAM
	}

	if _, err := global.GvaDockerCli.NetworkCreate(ctx, req.Name, options); err != nil {
		return err
	}
	return nil
}

/*
	Delete
	删除网络
*/
func (ns NetworkService) Delete(ctx context.Context, req NetworkDeleteStruct) error {

	for _, name := range req.Name {
		if err := global.GvaDockerCli.NetworkRemove(ctx, name); err != nil {
			return err
		}

	}
	return nil

}
