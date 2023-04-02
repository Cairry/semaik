package volumes

import (
	"dockerapi/sdk"
	"errors"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"log"
)

/*
	VolumesList
	卷列表
*/
func VolumesList() []Volume {

	volumes, err := sdk.DockerClient.VolumeList(sdk.DockerCtx, filters.Args{})
	if err != nil {
		log.Fatal(err)
	}

	// 定义卷列表切片
	result := make([]Volume, 0, len(volumes.Volumes))

	for _, volume := range volumes.Volumes {
		// 读取卷列表信息
		tmp := Volume{
			CreatedAt: volume.CreatedAt,
			Driver:    volume.Driver,
			Name:      volume.Name,
			Scope:     volume.Scope,
		}

		// 将卷列表信息假若都卷列表切片中
		result = append(result, tmp)

	}

	return result
}

/*
	VolumeSearch
	搜索卷
*/
func VolumeSearch(volumeName string) error {

	for _, v := range VolumesList() {
		if v.Name == volumeName {
			return errors.New("无法创建卷, 它已存在")
		}
	}

	return nil

}

/*
	VolumeCreate
	创建卷
*/
func VolumeCreate(req VolumeCreateStruct) error {

	// 获取创建卷选项配置
	options := volume.VolumeCreateBody{
		Driver:     req.Driver,
		DriverOpts: req.DriverOpts,
		Labels:     req.Labels,
		Name:       req.Name,
	}

	if err := VolumeSearch(req.Name); err != nil {
		return err
	}

	_, err := sdk.DockerClient.VolumeCreate(sdk.DockerCtx, options)

	if err != nil {
		return err
	}

	return nil

}

/*
	VolumeDelete
	删除卷
*/
func VolumeDelete(req VolumeDeleteStruct) error {

	for _, name := range req.Name {
		err := sdk.DockerClient.VolumeRemove(sdk.DockerCtx, name, true)
		if err != nil {
			return err
		}
	}
	return nil

}
