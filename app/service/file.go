package service

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"znews/app/dao"
	"znews/app/model"

	"github.com/gin-gonic/gin"
)

func Uploads(c *gin.Context, caseId string) error {
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
		err = c.SaveUploadedFile(file, path+"/"+caseId+"/"+file.Filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func Upload(c *gin.Context) error {
	account, _ := c.Get("account")

	// 解析上傳的檔案
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	tx := dao.GormSession.Begin()
	work := model.SohoWork{
		Account:  fmt.Sprintf("%v", account),
		FileName: file.Filename,
	}

	if err := tx.Model(&model.SohoWork{}).Create(&work).Error; err != nil {
		tx.Rollback()
		return err
	}

	path := os.Getenv("SOHO_WORK_PATH")

	if file.Size > 2<<20 {
		tx.Rollback()
		return errors.New(file.Filename + " 檔案大小超過 2 MB")
	}

	// 將文件保存到服務器上
	err = c.SaveUploadedFile(file, path+"/"+fmt.Sprintf("%v", account)+"/"+file.Filename)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
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

func SohoDownload(filename string, account string) string {
	lowerFilename := strings.ToLower(filename)

	if strings.HasSuffix(lowerFilename, ".doc") || strings.HasSuffix(lowerFilename, ".pdf") || strings.HasSuffix(lowerFilename, ".ppt") || strings.HasSuffix(lowerFilename, ".jpf") ||
		strings.HasSuffix(lowerFilename, ".png") || strings.HasSuffix(lowerFilename, ".txt") || strings.HasSuffix(lowerFilename, ".gif") {
		path := os.Getenv("SOHO_WORK_PATH")
		return path + "/" + account + "/" + filename
	} else {
		return ""
	}
}

func GetSohoWork(account string) ([]model.SohoWork, error) {
	work := []model.SohoWork{}

	err := dao.GormSession.Select("*").Where("account = ?", account).Find(&work).Error
	if err != nil {
		return nil, err
	} else {
		return work, nil
	}

}
