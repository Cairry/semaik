package common

import (
	"dockerapi/app/service/docker/node"
	"dockerapi/app/service/user"
	"dockerapi/global"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func init() {

	db, err := gorm.Open(sqlite.Open("./data/sql.db"), &gorm.Config{})
	if err != nil {
		log.Println("[error]: failed to connect database")
	}

	// 迁移 User 表
	_ = db.AutoMigrate(&user.User{})
	_ = db.AutoMigrate(&node.DockerNode{})

	global.GvaDatabase = db

}
