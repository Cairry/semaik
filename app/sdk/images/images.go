package images

import (
	"dockerapi/app/sdk"
	"github.com/docker/docker/api/types"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

/*
	ImagesList
	读取所有镜像
*/
func ImagesList() []ImageStruct {

	// 读取镜像列表
	images, err := sdk.DockerClient.ImageList(sdk.DockerCtx, types.ImageListOptions{All: true})
	if err != nil {
		log.Fatal(err)
	}

	// 初始化镜像列表，长度为读取到的数据列表长度
	result := make([]ImageStruct, 0, len(images))

	// 循环读取
	for _, image := range images {

		// 临时读取镜像列表
		tmp := ImageStruct{
			Created:    time.Unix(image.Created, 0),
			ID:         strings.Split(image.ID, ":")[1][:10],
			Repository: strings.Split(image.RepoTags[0], ":")[0],
			RepoTags:   strings.Split(image.RepoTags[0], ":")[1],
			Size:       image.Size,
		}

		// 将临时读取到的列表数据加入到初始化后的镜像列表中
		result = append(result, tmp)

	}

	// 打印镜像列表
	return result
}

/*
	PullImage
	拉取镜像
*/
func PullImage(imageName string) error {

	// 拉取镜像
	reader, err := sdk.DockerClient.ImagePull(sdk.DockerCtx, imageName, types.ImagePullOptions{})
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
	ImageDelete
	删除镜像
*/
func ImageDelete(req ImageOperationStruct) error {

	_, err := sdk.DockerClient.ImageRemove(sdk.DockerCtx, req.Name, types.ImageRemoveOptions{})

	return err

}
