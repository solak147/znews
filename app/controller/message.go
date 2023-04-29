package controller

import (
	"fmt"
	"net/http"
	"znews/app/model"
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
			"msg":  "GetMsgRecord failed : " + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Success",
			"data": data,
		})
	}

}

// @Summary GetMsgRecordDetail
// @Tags message
// @version 1.0
// @produce application/json
// @Security BearerAuth
// @param id path string true "訊息細節"
// @Success 200 string json successful return data
// @Router /message/{toAccount} [get]
func (m MsgsController) GetMsgRecordDetail(c *gin.Context) {

	account, _ := c.Get("account")
	to := c.Params.ByName("toAccount")

	data, err := service.GetMsgRecordDetail(fmt.Sprintf("%v", account), to)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  " GetMsgRecordDetail failed : " + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Success",
			"data": data,
		})
	}

}

// @Summary SendMsg
// @Tags message
// @version 1.0
// @produce application/json
// @Security BearerAuth
// @param msg body model.MsgSend true "傳送訊息"
// @Success 200 string json successful return data
// @Router /message/send [post]
func (m MsgsController) SendMsg(c *gin.Context) {

	account, _ := c.Get("account")

	var form model.MsgSend
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Form bind error": err.Error()})
		return
	}

	err := service.SendMsg(fmt.Sprintf("%v", account), form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  " SendMsg failed : " + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Success",
		})
	}

}
