package user

import (
	"crypto/md5"
	"dockerapi/app/response"
	"dockerapi/app/service"
	"dockerapi/app/service/user"
	"dockerapi/global"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRouter struct{}

/*
	Login
	用户登陆
*/
func (u UserRouter) Login(c *gin.Context) {

	var (
		use user.User
		req user.User
	)
	_ = c.ShouldBindJSON(&req)

	// 校验 Password
	arr := md5.Sum([]byte(req.Password))
	hashPassword := hex.EncodeToString(arr[:])

	// 查询用户信息
	err := global.GvaDatabase.Where("name = ?", req.Name).First(&use).Error
	if err == gorm.ErrRecordNotFound || hashPassword != use.Password {
		response.Fail(c, "用户不存在或密码错误", "登陆失败")
		return
	}

	tokenData, err, _ := service.JwtService.CreateToken(service.AppGuardName, req)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, tokenData, "登陆成功")

}

/*
	Register
	用户注册
*/
func (u UserRouter) Register(c *gin.Context) {

	var req user.User
	var use user.User

	_ = c.ShouldBindJSON(&req)
	global.GvaDatabase.Where("name = ?", req.Name).First(&use)
	if use.Name != "" {
		response.Fail(c, "用户已存在", "注册失败")
		return
	}

	arr := md5.Sum([]byte(req.Password))
	hashPassword := hex.EncodeToString(arr[:])
	global.GvaDatabase.Create(&user.User{
		Name:     req.Name,
		Password: hashPassword,
	})
	response.Success(c, nil, "注册成功")
}
