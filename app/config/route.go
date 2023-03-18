package config

import (
	"znews/app/controller"
	"znews/app/middleware"

	"github.com/gin-gonic/gin"
)

func RouteUsers(r *gin.Engine) {
	posts := r.Group("/v1/users")
	{
		posts.POST("/", controller.UserController().CreateUser)
		posts.GET("/", controller.UserController().GetUser)
	}

	login := r.Group("/login")
	{
		login.POST("/", controller.UserController().AuthHandler)
		login.GET("/:id", middleware.JWTAuthMiddleware(), controller.UserController().GetUser)
	}
}

// router.GET("/user/:name", func(c *gin.Context) {
// 	fmt.Println(c.FullPath())  // /user/:name/
// 	name := c.Param("name")
// 	c.String(http.StatusOK, "Hello %s", name)
// })
