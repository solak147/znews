package config

import (
	"znews/app/controller"

	"github.com/gin-gonic/gin"
)

func RouteUsers(r *gin.Engine) {
	posts := r.Group("/v1/users")
	{
		posts.POST("/", controller.NewUsersController().CreateUser)
		posts.GET("/", controller.QueryUsersController().GetUser)
	}
}
