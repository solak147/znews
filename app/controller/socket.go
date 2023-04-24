package controller

import (
	"znews/app/service"

	"github.com/gin-gonic/gin"
)

type SocketsController struct{}

func SocketController() SocketsController {
	return SocketsController{}
}

func (s SocketsController) Socket(c *gin.Context) {
	service.Socket(c)
}
