package docker

import (
	"dockerapi/app/response"
	"dockerapi/app/service/docker/node"
	"github.com/gin-gonic/gin"
)

type DockerNodeApi struct {
	dockerNode node.DockerNodeService
}

//	@Summary	创建 Docker 节点
//	@Tags		节点
//	@Produce	json
//	@Success 	200 {string}	json "{"code":200,"data":{},"msg":"创建成功"}"
//  @Failure	400 {string}	json "{"code":400,"data":{},"msg":"创建失败"}"
//	@Router		/api/v1/clouds/node/create [post]
func (d DockerNodeApi) Create(ctx *gin.Context) {

	var req node.DockerNode
	_ = ctx.ShouldBindJSON(&req)
	info, err := d.dockerNode.Create(req)
	if err != nil {
		response.Fail(ctx, err, "创建失败")
		return
	}
	response.Success(ctx, info, "创建成功")

}

//	@Summary	更新 Docker 节点
//	@Tags		节点
//	@Produce	json
//	@Success 	200 {string}	json "{"code":200,"data":{},"msg":"更新成功"}"
//  @Failure	400 {string}	json "{"code":400,"data":{},"msg":"更新失败"}"
//	@Router		/api/v1/clouds/node/update [put]
func (d DockerNodeApi) Update(ctx *gin.Context) {

	var req node.DockerNode
	_ = ctx.ShouldBindJSON(&req)
	info, err := d.dockerNode.Update(req)
	if err != nil {
		response.Fail(ctx, err, "更新失败")
		return
	}
	response.Success(ctx, info, "更新成功")
}

//	@Summary	删除 Docker 节点
//	@Tags		节点
//	@Produce	json
//	@Success 	200 {string}	json "{"code":200,"data":{},"msg":"删除成功"}"
//  @Failure	400 {string}	json "{"code":400,"data":{},"msg":"删除失败"}"
//	@Router		/api/v1/clouds/node/delete [delete]
func (d DockerNodeApi) Delete(ctx *gin.Context) {

	var req node.DockerNode
	_ = ctx.ShouldBindJSON(&req)
	info, err := d.dockerNode.Delete(req)
	if err != nil {
		response.Fail(ctx, err, "删除失败")
		return
	}
	response.Success(ctx, info, "删除成功")

}
