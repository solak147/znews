package main

import (
	"os"
	"znews/app/config"
	"znews/app/dao"
	"znews/app/model"

	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

// @title Gin swagger
// @version 1.8.10
// @description Gin swagger

// @contact.name kevin

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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
	migrateErr := db.AutoMigrate(&model.User{}, &model.Casem{}, &model.CaseFile{}, &model.SerialNo{})
	if migrateErr != nil {
		return
	}

	server := gin.Default()

	redisStore := persist.NewRedisStore(redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "redis:6379",
		DB:      0,
	}))

	config.CustomRouter(server, redisStore)

	err := server.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
