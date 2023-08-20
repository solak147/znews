package controller

import (
	"net/http"
	"znews/app/middleware"
	"znews/app/model"
	"znews/app/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PaysController struct{}

func PayController() PaysController {
	return PaysController{}
}

// @Summary 信用卡一次付清
// @Tags pay
// @version 1.0
// @produce application/json
// @param product formData string true "信用卡一次付清"
// @Success 200 boolean successful return boolean
// @Router /pay/creditAll [post]
func (p PaysController) CreditAll(c *gin.Context) {
	// 綠界金流的 API 位置
	apiUrl := "https://payment-stage.ecpay.com.tw/Cashier/AioCheckOut/V5"
	
	var form model.ProductForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	param, err := service.CreditAll(form.Card)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "Upload fail : " + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Success",
			"param": param,
			"api": apiUrl,
		},
	)}
}

// @Summary 付款結果
// @Tags pay
// @version 1.0
// @produce application/json
// @param pay path string true "付款結果"
// @Success 200 boolean successful return boolean
// @Router /result [post]
func (p PaysController) Result(c *gin.Context) {
	err := service.Result(c)

	if err != nil {
		middleware.Logger().WithFields(logrus.Fields{
			"title": "receive pay result error",
		}).Error(err.Error)
		c.String(http.StatusBadRequest, err.Error())
	} else {
		c.String(http.StatusOK, "1|OK")
	}	

}