package main

import (
	_ "dockerapi/common"
	"dockerapi/global"
	_ "dockerapi/routers/router"
	"log"
)

func main() {

	log.Println("Successfully.")

	_ = global.GvaGinEngine.Run()

}
