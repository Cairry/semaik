package networks

import (
	"dockerapi/sdk"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"log"
)

/*
	NetworksList
	网络列表
*/
func NetworksList() []Networks {

	networks, err := sdk.DockerClient.NetworkList(sdk.DockerCtx, types.NetworkListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// 初始化网络列表切片
	result := make([]Networks, 0, len(networks))

	// 循环读取
	for _, network := range networks {
		// 获取网络信息
		tmp := Networks{
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
	NetworkSearch
	搜索网络
*/
func NetworkSearch(networkName string) error {

	for _, v := range NetworksList() {
		if v.Name == networkName {
			return errors.New("无法创建网络, 它已存在")
		}
	}

	return nil
}

/*
	NetworkCreate
	创建网络
*/
func NetworkCreate(req NetworksCreateStruct) error {

	// 先判断是否存在
	if err := NetworkSearch(req.Name); err != nil {
		return err
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

	if _, err := sdk.DockerClient.NetworkCreate(sdk.DockerCtx, req.Name, options); err != nil {
		return err
	}
	return nil
}

/*
	NetworkDelete
	删除网络
*/
func NetworkDelete(req NetworkDeleteStruct) error {

	for _, name := range req.Name {
		if err := sdk.DockerClient.NetworkRemove(sdk.DockerCtx, name); err != nil {
			return err
		}

	}
	return nil

}
