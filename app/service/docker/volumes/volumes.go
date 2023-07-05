package volumes

import (
	"context"
	"dockerapi/global"
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
func (vs VolumeService) List(ctx context.Context) []VolumeList {

	volumes, err := global.GvaDockerCli.VolumeList(ctx, filters.Args{})
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
func (vs VolumeService) Search(ctx context.Context, info string) []VolumeList {

	volumes, err := global.GvaDockerCli.VolumeList(ctx, filters.Args{})
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
func (vs VolumeService) Create(ctx context.Context, req VolumeCreateStruct) error {

	for _, v := range vs.List(ctx) {
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

	_, err := global.GvaDockerCli.VolumeCreate(ctx, options)

	if err != nil {
		return err
	}

	return nil

}

/*
	Delete
	删除卷
*/
func (vs VolumeService) Delete(ctx context.Context, req VolumeDeleteStruct) error {

	for _, name := range req.Name {
		err := global.GvaDockerCli.VolumeRemove(ctx, name, true)
		if err != nil {
			return err
		}
	}
	return nil

}
