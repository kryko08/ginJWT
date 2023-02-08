package main

import (
	"GoProject/db"
	"github.com/gin-gonic/gin"
	"log"
)

func RouterSetUp() *gin.Engine {
	r := gin.Default()

	db.ConnectDB()

	return r
}

func main() {
	router := RouterSetUp()
	// run server
	err := router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		log.Fatal("Cannot run router", err)
	}
}
