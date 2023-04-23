package controller

import (
	"net/http"
	"znews/app/service"

	"github.com/gin-gonic/gin"
)

type FilesController struct{}

func FileController() FilesController {
	return FilesController{}
}

// @Summary 上傳檔案
// @Tags file
// @version 1.0
// @produce application/json
// @param files formData []string true "上傳檔案"
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

// @Summary 下載檔案
// @Tags file
// @version 1.0
// @produce application/json
// @param files path string true "下載檔案"
// @Success 200 file successful return file
// @Router /download/{filename} [get]
func (f FilesController) Download(c *gin.Context) {
	filename := c.Param("filename")

	if filepath := service.Download(filename); filepath == "" {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.File(filepath)
	}
}
