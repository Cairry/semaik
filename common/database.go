package common

import (
	"dockerapi/app/model"
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
	_ = db.AutoMigrate(&model.User{})

	global.GvaDatabase = db

}
