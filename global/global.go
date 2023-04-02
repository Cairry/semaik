package global

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (

	// GvaGinEngine 实例化 Gin 引擎
	GvaGinEngine = gin.Default()
	// GvaDatabase 实例化数据库
	GvaDatabase *gorm.DB
)
