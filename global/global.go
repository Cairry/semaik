package global

import (
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	GvaGinEngine = gin.Default()

	GvaDatabase *gorm.DB

	GvaDockerCli *client.Client
)
