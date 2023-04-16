package controller

import (
	"net/http"
	"znews/app/middleware"
	"znews/app/model"
	"znews/app/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CasesController struct{}

func CaseController() CasesController {
	return CasesController{}
}

// @Summary 發案/案件建立
// @Tags case
// @version 1.0
// @produce application/json
// @param case body model.CreateCase true "發案"
// @Success 200 boolean successful return boolean
// @Router /case/create [post]
func (ca CasesController) CreateCase(c *gin.Context) {

	var form model.CreateCase
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "表單綁定失敗 : " + err.Error(),
		})
		return
	}

	err := service.CreateCase(form)
	if err != nil {
		middleware.Logger().WithFields(logrus.Fields{
			"title":  "Create case Failed:",
			"accout": form.Account,
		}).Error(err.Error)

		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "發案失敗 : " + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Success",
		})
	}

}
