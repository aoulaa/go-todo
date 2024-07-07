package main

import (
	"github.com/gin-gonic/gin"
	"todo/api/router"
	"todo/internal/db"
)

func main() {
	route := gin.Default()

	// gin.SetMode(gin.ReleaseMode) //optional to not get warning
	// route.SetTrustedProxies([]string{"192.168.1.2"}) //to trust only a specific value

	db.ConnectDatabase()

	defer db.CloseDatabase()

	router.GetRoute(route)

	route.Run(":8080")
}
