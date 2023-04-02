package docker

import (
	"dockerapi/app/response"
	"dockerapi/sdk/containers"
	"github.com/gin-gonic/gin"
)

func ContainersList(c *gin.Context) {

	containersList := containers.ContainersList()
	response.Success(c, containersList, "请求成功")

}

func ContainerCreate(c *gin.Context) {

	var req containers.ContainerCreateStruct
	_ = c.ShouldBindJSON(&req)
	err := containers.ContainerCreate(req)
	if err != nil {
		response.Fail(c, nil, "创建失败")
	} else {
		response.Success(c, req, "创建成功")
	}

}

func ContainerOptions(c *gin.Context) {

	var req containers.ContainerOperationStruct
	_ = c.ShouldBindJSON(&req)
	err := containers.ContainerOptions(req)
	if err != nil {
		response.Fail(c, req.Operation, err.Error())
	} else {
		response.Success(c, req.Operation, "操作成功")
	}

}
