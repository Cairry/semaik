package volumes

import (
	"dockerapi/app/sdk"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"log"
	"strings"
)

type VolumeService struct{}

/*
	List
	卷列表
*/
func (vs VolumeService) List() []VolumeList {

	volumes, err := sdk.DockerClient.VolumeList(sdk.DockerCtx, filters.Args{})
	if err != nil {
		log.Fatal(err)
	}

	// 定义卷列表切片
	result := make([]VolumeList, 0, len(volumes.Volumes))

	for _, volume := range volumes.Volumes {
		// 读取卷列表信息
		tmp := VolumeList{
			CreatedAt: volume.CreatedAt,
			Driver:    volume.Driver,
			Name:      volume.Name,
			Scope:     volume.Scope,
		}

		// 将卷列表信息加入都卷列表切片中
		result = append(result, tmp)

	}

	return result
}

/*
	Search
	搜索卷
*/
func (vs VolumeService) Search(info string) []VolumeList {

	volumes, err := sdk.DockerClient.VolumeList(sdk.DockerCtx, filters.Args{})
	if err != nil {
		log.Fatal(err)
	}

	var volumesList []*types.Volume

	if len(info) != 0 {
		length, count := len(volumes.Volumes), 0
		for count < length {
			if strings.Contains(volumes.Volumes[count].Name, info) {
				volumesList = append(volumesList, volumes.Volumes[count])
				break
			}
			count++
		}
	} else {
		volumesList = volumes.Volumes
	}

	result := make([]VolumeList, 0, len(volumes.Volumes))

	for _, volume := range volumesList {

		tmp := VolumeList{
			CreatedAt: volume.CreatedAt,
			Driver:    volume.Driver,
			Name:      volume.Name,
			Scope:     volume.Scope,
			Status:    volume.Status,
		}

		result = append(result, tmp)

	}

	return result

}

/*
	Create
	创建卷
*/
func (vs VolumeService) Create(req VolumeCreateStruct) error {

	for _, v := range vs.List() {
		if v.Name == req.Name {
			return errors.New("无法创建卷, 它已存在")
		}
	}

	// 获取创建卷选项配置
	options := volume.VolumeCreateBody{
		Driver:     req.Driver,
		DriverOpts: req.DriverOpts,
		Labels:     req.Labels,
		Name:       req.Name,
	}

	_, err := sdk.DockerClient.VolumeCreate(sdk.DockerCtx, options)

	if err != nil {
		return err
	}

	return nil

}

/*
	Delete
	删除卷
*/
func (vs VolumeService) Delete(req VolumeDeleteStruct) error {

	for _, name := range req.Name {
		err := sdk.DockerClient.VolumeRemove(sdk.DockerCtx, name, true)
		if err != nil {
			return err
		}
	}
	return nil

}
