package docker

import (
	"dockerapi/app/response"
	"dockerapi/sdk/volumes"
	"github.com/gin-gonic/gin"
)

func VolumesList(c *gin.Context) {

	volumesList := volumes.VolumesList()
	response.Success(c, volumesList, "请求成功")

}

func VolumesCreate(c *gin.Context) {

	var req volumes.VolumeCreateStruct
	_ = c.ShouldBindJSON(&req)
	err := volumes.VolumeCreate(req)
	if err != nil {
		response.Fail(c, err.Error(), "创建失败")
	} else {
		response.Success(c, req, "创建成功")
	}

}

func VolumeDelete(c *gin.Context) {

	var req volumes.VolumeDeleteStruct
	_ = c.ShouldBindJSON(&req)
	err := volumes.VolumeDelete(req)
	if err != nil {
		response.Fail(c, err.Error(), "删除失败")
	} else {
		response.Success(c, req, "删除成功")
	}

}
