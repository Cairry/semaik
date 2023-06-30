package docker

import (
	"dockerapi/app/response"
	networks2 "dockerapi/app/sdk/networks"
	"github.com/gin-gonic/gin"
)

type NetworksApi struct {
	networkService networks2.NetworkService
}

//	@Summary	获取所有网络列表
//	@Tags		网络
//	@Produce	json
//	@Success 	200 {string}	json "{"code":200,"data":{},"msg":"请求成功"}"
//  @Failure	400 {string}	json "{"code":400,"data":{},"msg":"请求失败"}"
//	@Router		/api/v1/networks/list [get]
func (network NetworksApi) List(ctx *gin.Context) {

	networkList := network.networkService.List()
	response.Success(ctx, networkList, "请求成功")

}

//	@Summary	搜索网络
//	@Tags		网络
//	@Produce	json
//	@Success 	200 {string}	json "{"code":200,"data":{},"msg":"请求成功"}"
//  @Failure	400 {string}	json "{"code":400,"data":{},"msg":"请求失败"}"
//	@Router		/api/v1/networks/search [post]
func (network NetworksApi) Search(ctx *gin.Context) {

	var req networks2.NetworkList
	_ = ctx.ShouldBindJSON(&req)
	networkList := network.networkService.Search(req.Name)
	response.Success(ctx, networkList, "查询成功")
}

//	@Summary	创建网络
//	@Tags		网络
//	@Produce	json
//	@Success 	200 {string}	json "{"code":200,"data":{},"msg":"请求成功"}"
//  @Failure	400 {string}	json "{"code":400,"data":{},"msg":"请求失败"}"
//	@Router		/api/v1/networks/create [post]
func (network NetworksApi) Create(ctx *gin.Context) {

	var req networks2.NetworkCreateStruct
	_ = ctx.ShouldBindJSON(&req)
	err := network.networkService.Create(req)
	if err != nil {
		response.Fail(ctx, err.Error(), "创建失败")
	} else {
		response.Success(ctx, req, "创建成功")
	}

}

//	@Summary	删除网络
//	@Tags		网络
//	@Produce	json
//	@Success 	200 {string}	json "{"code":200,"data":{},"msg":"请求成功"}"
//  @Failure	400 {string}	json "{"code":400,"data":{},"msg":"请求失败"}"
//	@Router		/api/v1/networks/delete [post]
func (network NetworksApi) Delete(ctx *gin.Context) {

	var req networks2.NetworkDeleteStruct
	_ = ctx.ShouldBindJSON(&req)
	err := network.networkService.Delete(req)
	if err != nil {
		response.Fail(ctx, err.Error(), "删除失败")
	} else {
		response.Success(ctx, req, "删除成功")
	}

}
