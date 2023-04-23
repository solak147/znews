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

// @Summary 接案查詢
// @Tags case
// @version 1.0
// @produce application/json
// @param case path string true "接案查詢"
// @Success 200 json successful return json
// @Router /case/getAll [get]
func (ca CasesController) GetCase(c *gin.Context) {
	data, err, cnt := service.GetCase(c)

	if err != nil {
		middleware.Logger().WithFields(logrus.Fields{
			"title": "Get case failed:",
		}).Error(err.Error)

		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get case Failed :" + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Success",
			"data": data,
			"cnt":  cnt,
		})
	}

}

// @Summary 案件詳細資料
// @Tags case
// @version 1.0
// @produce application/json
// @param case path string true "案件詳細資料"
// @Success 200 json successful return json
// @Router /case/getDetail/{caseid} [get]
func (ca CasesController) GetCaseDetail(c *gin.Context) {
	caseId := c.Params.ByName("caseid")

	var (
		data  *model.Casem
		files []model.CaseFile
		err   error
	)

	if _, exit := c.Get("account"); exit {
		data, files, err = service.GetCaseDetail(caseId, true)
	} else {
		data, files, err = service.GetCaseDetail(caseId, false)
	}

	if err != nil {
		middleware.Logger().WithFields(logrus.Fields{
			"title": "Get case deatil failed:",
		}).Error(err.Error)

		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get case deatil failed :" + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"msg":   "Success",
			"data":  data,
			"files": files,
		})
	}
}
