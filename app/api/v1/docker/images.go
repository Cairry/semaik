package docker

import (
	"dockerapi/app/response"
	images2 "dockerapi/app/sdk/images"
	"github.com/gin-gonic/gin"
)

func ImagesList(c *gin.Context) {

	imageInfo := images2.ImagesList()
	response.Success(c, imageInfo, "请求成功")

}

func ImageDelete(c *gin.Context) {

	var req images2.ImageOperationStruct
	_ = c.ShouldBindJSON(&req)
	if err := images2.ImageDelete(req); err != nil {
		response.Fail(c, err.Error(), "删除失败")
	} else {
		response.Success(c, req, "删除成功")
	}

}

func ImagePull(c *gin.Context) {

	var req images2.ImageOperationStruct
	_ = c.ShouldBindJSON(&req)
	err := images2.PullImage(req.Name)
	if err != nil {
		response.Fail(c, err.Error(), "拉取失败")
	} else {
		response.Success(c, req.Name, "拉取成功")
	}

}
