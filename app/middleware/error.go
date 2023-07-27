package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ErrorMiddleware(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			Logger().WithFields(logrus.Fields{
				"title": "非預期錯誤:",
			}).Error(err)

			c.JSON(http.StatusBadRequest, gin.H{
				"code": -3,
				"msg":  "非預期錯誤",
			})

			c.Abort()
		}
	}()
	c.Next()
}
