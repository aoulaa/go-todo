package router

import (
	"github.com/gin-gonic/gin"
	"todo/api/middleware"
	"todo/api/rest"
)

func GetRoute(r *gin.Engine) {
	r.POST("/api/signup", rest.Signup)
	r.POST("/api/login", rest.Login)

	r.Use(middleware.RequireAuth)

	r.POST("/api/logout", rest.Logout)

	userRouter := r.Group("/api/users")
	{
		userRouter.GET("/", rest.GetUsers)
		userRouter.GET("/:id", rest.GetUser)
		userRouter.POST("/", rest.AddUser)
		userRouter.PUT("/:id", rest.UpdateUser)
		userRouter.DELETE("/:id", rest.DeleteUser)
	}

}
