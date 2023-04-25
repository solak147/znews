package controller

import (
	"fmt"
	"net/http"
	"znews/app/service"

	"github.com/gin-gonic/gin"
)

type MsgsController struct{}

func MsgController() MsgsController {
	return MsgsController{}
}

func (m MsgsController) GetMsgRecord(c *gin.Context) {

	account, _ := c.Get("account")

	data, err := service.GetMsgRecord(fmt.Sprintf("%v", account))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "Upload fail : " + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Success",
			"data": data,
		})
	}

}
