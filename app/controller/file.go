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
