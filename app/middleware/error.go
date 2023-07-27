package middleware

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ErrorMiddleware(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			_, file, line, _ := runtime.Caller(2) // 獲取呼叫者的信息（1代表ErrorMiddleware本身，2代表ErrorMiddleware的呼叫者）
			Logger().WithFields(logrus.Fields{
				"name": "非預期錯誤:",
			}).Error("file:", file, " | line:", line, " | error:", err)

			c.JSON(http.StatusBadRequest, gin.H{
				"code": -3,
				"msg":  "非預期錯誤",
			})

			c.Abort()
		}
	}()
	c.Next()
}
