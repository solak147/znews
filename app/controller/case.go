package controller

import (
	"fmt"
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

	err := service.CreateCase(c)
	if err != nil {
		middleware.Logger().WithFields(logrus.Fields{
			"title": "Create case Failed:",
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
	caseId := c.Params.ByName("caseId")

	var (
		data  *model.Casem
		files []model.CaseFile
		isVip bool
		err   error
	)

	if val, exit := c.Get("account"); exit {
		data, files, isVip, err = service.GetCaseDetail(caseId, fmt.Sprintf("%v", val))
	} else {
		data, files, isVip, err = service.GetCaseDetail(caseId, fmt.Sprintf("%v", val))
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
			"isVip": isVip,
		})
	}
}

// @Summary 報價
// @Tags case
// @version 1.0
// @produce application/json
// @Security BearerAuth
// @param quote body string true "新增報價"
// @Success 200 string json successful return data
// @Router /case/Quote [post]
func (ca CasesController) Quote(c *gin.Context) {
	account, _ := c.Get("account")

	var form model.QuoteForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Form bind error": err.Error()})
		return
	}

	err := service.Quote(fmt.Sprintf("%v", account), form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Success",
		})
	}
}

// @Summary 報價紀錄
// @Tags case
// @version 1.0
// @produce application/json
// @Security BearerAuth
// @param quote body string true "報價紀錄"
// @Success 200 string json successful return data
// @Router /case/quoteRecord [get]
func (ca CasesController) QuoteRecord(c *gin.Context) {
	account, _ := c.Get("account")
	deal := c.Params.ByName("deal")

	data, err := service.QuoteRecord(fmt.Sprintf("%v", account), deal)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Success",
			"data": data,
		})
	}
}

// @Summary 報價前檢查
// @Tags case
// @version 1.0
// @produce application/json
// @Security BearerAuth
// @param quote body string true "報價前檢查"
// @Success 200 string json successful return data
// @Router /case/chkBefQuote [get]
func (ca CasesController) ChkBefQuote(c *gin.Context) {
	account, _ := c.Get("account")
	caseId := c.Params.ByName("caseId")

	err := service.ChkBefQuote(fmt.Sprintf("%v", account), caseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Success",
		})
	}
}

// @Summary 取得案件流程
// @Tags case
// @version 1.0
// @produce application/json
// @Security BearerAuth
// @param flow path string true "取得案件流程"
// @Success 200 string json successful return data
// @Router /case/getflow [get]
func (ca CasesController) GetFlow(c *gin.Context) {
	caseId := c.Params.ByName("caseId")

	casem, flow, err := service.GetFlow(caseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Success",
			"case": casem,
			"flow": flow,
		})
	}
}
