package main

import (
	"os"
	"znews/app/config"
	"znews/app/dao"
	"znews/app/model"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		panic(envErr)
	}

	port := os.Getenv("PORT")
	dbConfig := os.Getenv("DB_CONFIG")
	db, ormErr := dao.Initialize(dbConfig)
	if ormErr != nil {
		panic(ormErr)
	}
	migrateErr := db.AutoMigrate(&model.User{})
	if migrateErr != nil {
		return
	}

	server := gin.Default()
	server.GET("/hc", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "health check",
		})
	})
	config.RouteUsers(server)
	err := server.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
