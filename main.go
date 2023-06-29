package main

import (
	_ "dockerapi/common"
	"dockerapi/global"
	_ "dockerapi/routers"
	"log"
)

func main() {

	log.Println("Successfully.")

	_ = global.GvaGinEngine.Run(":8081")

}
