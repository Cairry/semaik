package docker

import (
	"dockerapi/app/response"
	containers2 "dockerapi/app/sdk/containers"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

func ContainersList(c *gin.Context) {

	containersList := containers2.ContainersList()
	response.Success(c, containersList, "请求成功")

}

func ContainerCreate(c *gin.Context) {

	var req containers2.ContainerCreateStruct
	_ = c.ShouldBindJSON(&req)
	err := containers2.ContainerCreate(req)
	if err != nil {
		response.Fail(c, nil, "创建失败")
	} else {
		response.Success(c, req, "创建成功")
	}

}

func ContainerOptions(c *gin.Context) {

	var req containers2.ContainerOperationStruct
	_ = c.ShouldBindJSON(&req)
	err := containers2.ContainerOptions(req)
	if err != nil {
		response.Fail(c, req.Operation, err.Error())
	} else {
		response.Success(c, req.Operation, "操作成功")
	}

}

func ContainerLogs(c *gin.Context) {
	var req containers2.ContainerLogsStruct
	_ = c.ShouldBindJSON(&req)
	out, err := containers2.ContainerLogs(req)
	if err != nil {
		response.Fail(c, out, err.Error())
	} else {
		response.Success(c, out, "操作成功")
	}

}

func ContainerState(c *gin.Context) {

	containerID, ok := c.Params.Get("id")
	if !ok {
		response.Fail(c, ok, "请求失败")
		return
	}

	res, err := containers2.ContainerState(containerID)
	if err != nil {
		response.Fail(c, err, "请求失败")
		return
	} else {
		response.Success(c, res, "请求成功")
		return
	}

}

func ContainerWsSsh(c *gin.Context) {

	// 升级连接为 WebSocket
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024 * 1024 * 10,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// websocket 握手
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		response.Fail(c, err.Error(), "请求失败")
		return
	}

	defer wsConn.Close()

	// 容器信息
	req := containers2.ContainerWsSshStruct{
		ContainerID: c.Query("containerID"),
		User:        c.Query("user"),
		Command:     []string{c.Query("command")},
	}

	err = containers2.ContainerSsh(req, c, wsConn)
	if err != nil {
		response.Fail(c, err, "连接失败")
		return
	}
	response.Success(c, "", "连接成功")

}
