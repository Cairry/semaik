package images

import (
	"context"
	"dockerapi/global"
	"dockerapi/utils"
	"github.com/docker/docker/api/types"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type ImageService struct{}

/*
	List
	读取所有镜像
*/
func (is ImageService) List(ctx context.Context) []ImageStruct {

	// 读取镜像列表
	images, err := global.GvaDockerCli.ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		log.Fatal(err)
	}

	// 初始化镜像列表，长度为读取到的数据列表长度
	result := make([]ImageStruct, 0, len(images))

	// 循环读取
	for _, image := range images {

		// 临时读取镜像列表
		tmp := ImageStruct{
			Created: time.Unix(image.Created, 0),
			ID:      strings.Split(image.ID, ":")[1][:10],
			Tags:    image.RepoTags,
			Size:    utils.FormatFileSize(image.Size),
		}

		// 将临时读取到的列表数据加入到初始化后的镜像列表中
		result = append(result, tmp)

	}

	// 打印镜像列表
	return result
}

/*
	Pull
	拉取镜像
*/
func (is ImageService) Pull(ctx context.Context, imageName string) error {

	// 拉取镜像
	reader, err := global.GvaDockerCli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		log.Println(err)
		return err
	}
	// 拉取结束后关闭io操作
	defer reader.Close()
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

/*
	Delete
	删除镜像
*/
func (is ImageService) Delete(ctx context.Context, req ImageOperationStruct) error {

	_, err := global.GvaDockerCli.ImageRemove(ctx, req.Name, types.ImageRemoveOptions{})

	return err

}

/*
	Search
	搜索镜像
*/
func (is ImageService) Search(ctx context.Context, info string) []ImageStruct {

	images, err := global.GvaDockerCli.ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		log.Fatal(err)
	}

	var newList []types.ImageSummary

	// 判断用户是否有搜索内容
	if len(info) != 0 {
		length, count := len(images), 0
		for count < length {
			for _, v := range images[count].RepoTags {
				if strings.Contains(v, info) {
					newList = append(newList, images[count])
					break
				}
			}
			count++
		}
	} else {
		newList = images
	}

	// 初始化镜像列表，长度为读取到的数据列表长度
	result := make([]ImageStruct, 0, len(newList))

	for _, image := range newList {
		tmp := ImageStruct{
			Created: time.Unix(image.Created, 0),
			ID:      image.ID,
			Tags:    image.RepoTags,
			Size:    utils.FormatFileSize(image.Size),
		}
		result = append(result, tmp)
	}

	return result

}
