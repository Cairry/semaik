package docker

import (
	"dockerapi/app/response"
	volumes2 "dockerapi/app/sdk/volumes"
	"github.com/gin-gonic/gin"
)

func VolumesList(c *gin.Context) {

	volumesList := volumes2.VolumesList()
	response.Success(c, volumesList, "请求成功")

}

func VolumesCreate(c *gin.Context) {

	var req volumes2.VolumeCreateStruct
	_ = c.ShouldBindJSON(&req)
	err := volumes2.VolumeCreate(req)
	if err != nil {
		response.Fail(c, err.Error(), "创建失败")
	} else {
		response.Success(c, req, "创建成功")
	}

}

func VolumeDelete(c *gin.Context) {

	var req volumes2.VolumeDeleteStruct
	_ = c.ShouldBindJSON(&req)
	err := volumes2.VolumeDelete(req)
	if err != nil {
		response.Fail(c, err.Error(), "删除失败")
	} else {
		response.Success(c, req, "删除成功")
	}

}
