package auth

import (
	"crypto/md5"
	"dockerapi/app/response"
	"dockerapi/app/service"
	"dockerapi/global"
	"dockerapi/model"
	"encoding/hex"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {

	var (
		user model.User
		req  model.User
	)
	_ = c.ShouldBindJSON(&req)
	// 查询用户信息
	global.GvaDatabase.First(&user, "Name = ?", req.Name)
	if req.Name != user.Name {
		response.Fail(c, nil, "用户不存在")
		return
	}

	// 校验 Password
	arr := md5.Sum([]byte(req.Password))
	hashPassword := hex.EncodeToString(arr[:])
	if hashPassword != user.Password {
		response.Fail(c, nil, "密码错误")
		return
	} else {
		tokenData, err, _ := service.JwtService.CreateToken(service.AppGuardName, req)
		if err != nil {
			response.Fail(c, nil, err.Error())
			return
		}
		response.Success(c, tokenData, "登陆成功")
	}

}

func Register(c *gin.Context) {

	var (
		user model.User
		req  model.User
	)

	_ = c.ShouldBindJSON(&req)
	global.GvaDatabase.First(&user, "Name = ?", req.Name)
	if req.Name == user.Name {
		response.Fail(c, nil, "用户已存在")
		return
	}

	arr := md5.Sum([]byte(req.Password))
	hashPassword := hex.EncodeToString(arr[:])
	global.GvaDatabase.Create(&model.User{
		Name:     req.Name,
		Password: hashPassword,
	})
	response.Success(c, nil, "注册成功")
}
