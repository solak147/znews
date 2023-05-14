package controller

import (
	"fmt"
	"net/http"
	"znews/app/service"

	"github.com/gin-gonic/gin"
)

type FilesController struct{}

func FileController() FilesController {
	return FilesController{}
}

// @Summary 上傳多個檔案
// @Tags file
// @version 1.0
// @produce application/json
// @param files formData []string true "上傳多個檔案"
// @Success 200 boolean successful return boolean
// @Router /uploads [post]
// func (f FilesController) Uploads(c *gin.Context) {
// 	err := service.Uploads(c)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"code": -1,
// 			"msg":  "Uploads fail : " + err.Error(),
// 		})
// 	} else {
// 		c.JSON(http.StatusOK, gin.H{
// 			"code": 0,
// 			"msg":  "Success",
// 		})
// 	}
// }

// @Summary 上傳單個檔案
// @Tags file
// @version 1.0
// @produce application/json
// @param file formData string true "上傳單個檔案"
// @Success 200 boolean successful return boolean
// @Router /upload [post]
func (f FilesController) Upload(c *gin.Context) {
	err := service.Upload(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "Upload fail : " + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Success",
		})
	}
}

// @Summary 作品資料
// @Tags file
// @version 1.0
// @produce application/json
// @param file path string true "作品資料"
// @Success 200 boolean successful return boolean
// @Router /sohowork [get]
func (f FilesController) GetSohoWork(c *gin.Context) {
	account, _ := c.Get("account")
	accountPath := c.Query("account")
	param := c.Param("param")

	var tmp string
	if accountPath == "" || accountPath == "x" {
		tmp = fmt.Sprintf("%v", account)
	} else {
		tmp = accountPath
	}

	data, err := service.GetSohoWork(tmp, param)
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

// @Summary 下載案件檔案
// @Tags file
// @version 1.0
// @produce application/json
// @param files path string true "下載案件檔案"
// @Success 200 file successful return file
// @Router /download/{filename} [get]
func (f FilesController) Download(c *gin.Context) {
	caseId := c.Param("caseId")
	filename := c.Param("filename")

	if filepath := service.Download(caseId, filename); filepath == "" {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.File(filepath)
	}
}

// @Summary 下載作品檔案
// @Tags file
// @version 1.0
// @produce application/json
// @param files path string true "下載作品檔案"
// @Success 200 file successful return file
// @Router /download/{filename} [get]
func (f FilesController) SohoDownload(c *gin.Context) {
	if filepath := service.SohoDownload(c); filepath == "" {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.File(filepath)
	}
}

// @Summary 刪除作品檔案
// @Tags file
// @version 1.0
// @produce application/json
// @param files path string true "刪除作品檔案"
// @Success 200 string json successful return data
// @Router /file/sohowork/{filename} [delete]
func (f FilesController) DeleteSohoWork(c *gin.Context) {

	if err := service.DeleteSohoWork(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "Delete fail : " + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Success",
		})
	}
}
