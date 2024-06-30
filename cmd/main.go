package main

import (
	"github.com/gin-gonic/gin"
	"todo/internal/db"
	"todo/internal/rest"
)

func main() {
	// gin.SetMode(gin.ReleaseMode) //optional to not get warning
	// route.SetTrustedProxies([]string{"192.168.1.2"}) //to trust only a specific value
	route := gin.Default()
	db.ConnectDatabase()
	route.POST("/user", rest.AddUser)
	route.Run(":8080")
}
