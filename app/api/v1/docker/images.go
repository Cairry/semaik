package docker

import (
	"dockerapi/app/response"
	images2 "dockerapi/app/service/docker/images"
	"dockerapi/middleware"
	"github.com/gin-gonic/gin"
)

type ImagesApi struct {
	imageService images2.ImageService
}

//	@Summary	获取所有镜像
//	@Tags		镜像
//	@Produce	json
//	@Success 	200 {string}	json "{"code":200,"data":{},"msg":"请求成功"}"
//  @Failure	400 {string}	json "{"code":400,"data":{},"msg":"请求失败"}"
//	@Router		/api/v1/clouds/node/:node_name/images/list [get]
func (image ImagesApi) List(ctx *gin.Context) {

	// docker node
	middleware.InitCli(ctx)

	imageInfo := image.imageService.List(ctx)
	response.Success(ctx, imageInfo, "请求成功")

}

//	@Summary	删除镜像
//	@Tags		镜像
//	@Produce	json
// 	@Param		name		path	string	true	"镜像名称"
//	@Success 	200 {string}	json "{"code":200,"data":{},"msg":"请求成功"}"
//  @Failure	400 {string}	json "{"code":400,"data":{},"msg":"请求失败"}"
//	@Router		/api/v1/clouds/node/:node_name/images/delete [post]
func (image ImagesApi) Delete(ctx *gin.Context) {

	// docker node
	middleware.InitCli(ctx)

	var req images2.ImageOperationStruct
	_ = ctx.ShouldBindJSON(&req)
	if err := image.imageService.Delete(ctx, req); err != nil {
		response.Fail(ctx, err.Error(), "删除失败")
		return
	} else {
		response.Success(ctx, req, "删除成功")
	}

}

//	@Summary	拉取镜像
//	@Tags		镜像
//	@Produce	json
// 	@Param		name		path	string	true	"镜像名称"
//	@Success 	200 {string}	json "{"code":200,"data":{},"msg":"请求成功"}"
//  @Failure	400 {string}	json "{"code":400,"data":{},"msg":"请求失败"}"
//	@Router		/api/v1/clouds/node/:node_name/images/pull [post]
func (image ImagesApi) Pull(ctx *gin.Context) {

	// docker node
	middleware.InitCli(ctx)

	var req images2.ImageOperationStruct
	_ = ctx.ShouldBindJSON(&req)
	err := image.imageService.Pull(ctx, req.Name)
	if err != nil {
		response.Fail(ctx, err.Error(), "拉取失败")
		return
	} else {
		response.Success(ctx, req.Name, "拉取成功")
	}

}

//	@Summary	搜索镜像
//	@Tags		镜像
//	@Produce	json
// 	@Param		name		path	string	true	"镜像名称"
//	@Success 	200 {string}	json "{"code":200,"data":{},"msg":"请求成功"}"
//  @Failure	400 {string}	json "{"code":400,"data":{},"msg":"请求失败"}"
//	@Router		/api/v1/clouds/node/:node_name/images/search [post]
func (image ImagesApi) Search(ctx *gin.Context) {

	// docker node
	middleware.InitCli(ctx)

	var req images2.ImageOperationStruct
	_ = ctx.ShouldBindJSON(&req)
	info := image.imageService.Search(ctx, req.Name)
	response.Success(ctx, info, "查询成功")

}
