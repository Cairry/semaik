package docker

import (
	"dockerapi/app/response"
	containers2 "dockerapi/app/sdk/containers"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type ContainersApi struct {
	containerService containers2.ContainerService
}

//	@Summary	获取所有容器列表
//	@Tags		容器
//	@Produce	json
//	@Success 	200 {string}	json "{"code":200,"data":{},"msg":"请求成功"}"
//  @Failure	400 {string}	json "{"code":400,"data":{},"msg":"请求失败"}"
//	@Router		/api/v1/containers/list [get]
func (container ContainersApi) List(ctx *gin.Context) {

	containersList := container.containerService.List()
	response.Success(ctx, containersList, "请求成功")

}

//	@Summary	获取指定容器
//	@Tags		容器
//	@Produce	json
//	@Success 	200 {string}	json "{"code":200,"data":{},"msg":"请求成功"}"
//  @Failure	400 {string}	json "{"code":400,"data":{},"msg":"请求失败"}"
//	@Router		/api/v1/containers/Search [post]
func (container ContainersApi) Search(ctx *gin.Context) {

	var req containers2.ContainerListStruct
	_ = ctx.ShouldBindJSON(&req)
	containersList := container.containerService.Search(req.Name)
	response.Success(ctx, containersList, "查询成功")

}

//	@Summary	创建容器
//	@Tags		容器
//	@Produce	json
//	@Success 	200 {string}	json "{"code":200,"data":{},"msg":"请求成功"}"
//  @Failure	400 {string}	json "{"code":400,"data":{},"msg":"请求失败"}"
//	@Router		/api/v1/containers/create [post]
func (container ContainersApi) Create(ctx *gin.Context) {

	var req containers2.ContainerCreateStruct
	_ = ctx.ShouldBindJSON(&req)
	err := container.containerService.Create(req)
	if err != nil {
		response.Fail(ctx, nil, "创建失败")
	} else {
		response.Success(ctx, req, "创建成功")
	}

}

//	@Summary	容器操作选项
//	@Tags		容器
//	@Produce	json
//	@Success 	200 {string}	json "{"code":200,"data":{},"msg":"请求成功"}"
//  @Failure	400 {string}	json "{"code":400,"data":{},"msg":"请求失败"}"
//	@Router		/api/v1/containers/option [post]
func (container ContainersApi) Options(ctx *gin.Context) {

	var req containers2.ContainerOperationStruct
	_ = ctx.ShouldBindJSON(&req)
	err := container.containerService.Options(req)
	if err != nil {
		response.Fail(ctx, req.Operation, err.Error())
	} else {
		response.Success(ctx, req.Operation, "操作成功")
	}

}

//	@Summary	容器日志
//	@Tags		容器
//	@Produce	json
//	@Param		containerName 	path	string	true	"容器名称"
//	@Param		containerId		path	string	true	"容器ID"
//	@Param		isWatch			path	bool	false	"是否监听日志"
//	@Param		mode			path	string	false	"获取日志范围, all,1m,10m..."
//	@Success 	200 {string}	json "{"code":200,"data":{},"msg":"请求成功"}"
//  @Failure	400 {string}	json "{"code":400,"data":{},"msg":"请求失败"}"
//	@Router		/api/v1/containers/log [post]
func (container ContainersApi) Logs(ctx *gin.Context) {
	var req containers2.ContainerLogsStruct
	_ = ctx.ShouldBindJSON(&req)
	out, err := container.containerService.Logs(req)
	if err != nil {
		response.Fail(ctx, out, err.Error())
	} else {
		response.Success(ctx, out, "操作成功")
	}

}

//	@Summary	获取容器资源
//	@Tags		容器
//	@Produce	json
//	@Param		containerId 	path	string	true	"容器ID"
//	@Success 	200 {string}	json "{"code":200,"data":{},"msg":"请求成功"}"
//  @Failure	400 {string}	json "{"code":400,"data":{},"msg":"请求失败"}"
//	@Router		/api/v1/containers/state/:id [get]
func (container ContainersApi) State(ctx *gin.Context) {

	containerID, ok := ctx.Params.Get("id")
	if !ok {
		response.Fail(ctx, ok, "请求失败")
		return
	}

	res, err := container.containerService.State(containerID)
	if err != nil {
		response.Fail(ctx, err, "请求失败")
		return
	} else {
		response.Success(ctx, res, "请求成功")
		return
	}

}

//	@Summary	容器终端
//	@Tags		容器
//	@Produce	json
//	@Param		containerId 	path	string	true	"容器ID"
//	@Param		user			path	string	true	"user"
//	@Param		command			path	string	true	"执行命令"
//	@Success 	200 {string}	json "{"code":200,"data":{},"msg":"请求成功"}"
//  @Failure	400 {string}	json "{"code":400,"data":{},"msg":"请求失败"}"
//	@Router		/api/v1/containers/state/:id [get]
func (container ContainersApi) Ssh(ctx *gin.Context) {

	// 升级连接为 WebSocket
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024 * 1024 * 10,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// websocket 握手
	wsConn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		response.Fail(ctx, err.Error(), "请求失败")
		return
	}

	defer wsConn.Close()

	// 容器信息
	req := containers2.ContainerWsSshStruct{
		ContainerID: ctx.Query("containerID"),
		User:        ctx.Query("user"),
		Command:     []string{ctx.Query("command")},
	}

	err = container.containerService.Ssh(req, ctx, wsConn)
	if err != nil {
		response.Fail(ctx, err, "连接失败")
		return
	}
	response.Success(ctx, "", "连接成功")

}
