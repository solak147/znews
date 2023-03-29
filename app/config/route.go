package config

import (
	"time"
	_ "znews/docs"

	"znews/app/controller"
	"znews/app/middleware"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func CustomRouter(r *gin.Engine, m *persist.RedisStore) {

	r.Use(corsMiddleware())

	posts := r.Group("/v1/users")
	{
		posts.POST("/", controller.UserController().CreateUser)
		posts.GET("/", controller.UserController().GetUser)
	}

	login := r.Group("/member")
	{
		login.POST("/login", controller.UserController().Login)
		login.GET("/:id", middleware.JWTAuthMiddleware(), cache.CacheByRequestURI(m, 2*time.Hour), controller.UserController().GetUser)
	}

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 1天

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// router.GET("/user/:name", func(c *gin.Context) {
// 	fmt.Println(c.FullPath())  // /user/:name/
// 	name := c.Param("name")
// 	c.String(http.StatusOK, "Hello %s", name)
// })

// server.GET("/hc", func(c *gin.Context) {
// 	c.JSON(200, gin.H{
// 		"message": "health check",
// 	})
// })
