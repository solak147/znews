package service

import (
	"os"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["files[]"]
	path := os.Getenv("CASE_FILE")

	for _, file := range files {

		// 將文件保存到服務器上
		err = c.SaveUploadedFile(file, path+"/"+file.Filename)
		if err != nil {
			return err
		}
	}
	return nil
}
