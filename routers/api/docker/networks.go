package docker

import (
	"dockerapi/app/response"
	"dockerapi/sdk/networks"
	"github.com/gin-gonic/gin"
)

func NetworksList(c *gin.Context) {

	networkList := networks.NetworksList()
	response.Success(c, networkList, "请求成功")

}

func NetworksCreate(c *gin.Context) {

	var req networks.NetworksCreateStruct
	_ = c.ShouldBindJSON(&req)
	err := networks.NetworkCreate(req)
	if err != nil {
		response.Fail(c, err.Error(), "创建失败")
	} else {
		response.Success(c, req, "创建成功")
	}

}

func NetworkDelete(c *gin.Context) {

	var req networks.NetworkDeleteStruct
	_ = c.ShouldBindJSON(&req)
	err := networks.NetworkDelete(req)
	if err != nil {
		response.Fail(c, err.Error(), "删除失败")
	} else {
		response.Success(c, req, "删除成功")
	}

}
