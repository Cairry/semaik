package docker

import (
	"dockerapi/app/response"
	volumes2 "dockerapi/app/sdk/volumes"
	"github.com/gin-gonic/gin"
)

type VolumesApi struct {
	volumeService volumes2.VolumeService
}

//	@Summary	获取所有卷
//	@Tags		卷
//	@Produce	json
//	@Success 	200 {string}	json "{"code":200,"data":{},"msg":"请求成功"}"
//  @Failure	400 {string}	json "{"code":400,"data":{},"msg":"请求失败"}"
//	@Router		/api/v1/volumes/list [get]
func (volume VolumesApi) List(ctx *gin.Context) {

	volumesList := volume.volumeService.List()
	response.Success(ctx, volumesList, "请求成功")

}

//	@Summary	搜索卷
//	@Tags		卷
//	@Produce	json
//	@Success 	200 {string}	json "{"code":200,"data":{},"msg":"请求成功"}"
//  @Failure	400 {string}	json "{"code":400,"data":{},"msg":"请求失败"}"
//	@Router		/api/v1/volumes/search [post]
func (volume VolumesApi) Search(ctx *gin.Context) {

	var req volumes2.VolumeList
	_ = ctx.ShouldBindJSON(&req)
	volumeList := volume.volumeService.Search(req.Name)
	response.Success(ctx, volumeList, "查询成功")

}

//	@Summary	创建卷
//	@Tags		卷
//	@Produce	json
//	@Success 	200 {string}	json "{"code":200,"data":{},"msg":"请求成功"}"
//  @Failure	400 {string}	json "{"code":400,"data":{},"msg":"请求失败"}"
//	@Router		/api/v1/volumes/create [post]
func (volume VolumesApi) Create(ctx *gin.Context) {

	var req volumes2.VolumeCreateStruct
	_ = ctx.ShouldBindJSON(&req)
	err := volume.volumeService.Create(req)
	if err != nil {
		response.Fail(ctx, err.Error(), "创建失败")
	} else {
		response.Success(ctx, req, "创建成功")
	}

}

//	@Summary	删除卷
//	@Tags		卷
//	@Produce	json
//	@Success 	200 {string}	json "{"code":200,"data":{},"msg":"请求成功"}"
//  @Failure	400 {string}	json "{"code":400,"data":{},"msg":"请求失败"}"
//	@Router		/api/v1/volumes/delete [post]
func (volume VolumesApi) Delete(ctx *gin.Context) {

	var req volumes2.VolumeDeleteStruct
	_ = ctx.ShouldBindJSON(&req)
	err := volume.volumeService.Delete(req)
	if err != nil {
		response.Fail(ctx, err.Error(), "删除失败")
	} else {
		response.Success(ctx, req, "删除成功")
	}

}
