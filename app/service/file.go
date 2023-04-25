package service

import (
	"errors"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["files[]"]
	path := os.Getenv("Case_FILE_PATH")

	for _, file := range files {

		if file.Size > 2<<20 {
			return errors.New(file.Filename + " 檔案大小超過 2 MB")
		}

		// 將文件保存到服務器上
		err = c.SaveUploadedFile(file, path+"/"+file.Filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func Download(filename string) string {
	lowerFilename := strings.ToLower(filename)

	if strings.HasSuffix(lowerFilename, ".doc") || strings.HasSuffix(lowerFilename, ".pdf") || strings.HasSuffix(lowerFilename, ".ppt") || strings.HasSuffix(lowerFilename, ".jpf") ||
		strings.HasSuffix(lowerFilename, ".png") || strings.HasSuffix(lowerFilename, ".txt") || strings.HasSuffix(lowerFilename, ".gif") {
		path := os.Getenv("Case_FILE_PATH")
		return path + "/" + filename
	} else {
		return ""
	}
}
