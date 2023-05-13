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
			err := os.RemoveAll(path + "/" + caseId)
			if err != nil {
				return err
			}
			return err
		}
	}
	return nil
}

func Upload(c *gin.Context) error {
	account, _ := c.Get("account")
	param := c.Param("param")

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

	// 頭像
	if param == "avatar" {
		work.FileType = "1"
	} else if param == "work" {
		work.FileType = "0"
	}

	if err := tx.Model(&model.SohoWork{}).Create(&work).Error; err != nil {
		tx.Rollback()
		return err
	}

	var path string
	if param == "avatar" {
		path = os.Getenv("SOHO_WORK_PATH") + "/" + fmt.Sprintf("%v", account) + "/avatar/" + file.Filename
	} else {
		path = os.Getenv("SOHO_WORK_PATH") + "/" + fmt.Sprintf("%v", account) + "/" + file.Filename
	}

	if file.Size > 2<<20 {
		tx.Rollback()
		return errors.New(file.Filename + " 檔案大小超過 2 MB")
	}

	// 將文件保存到服務器上
	err = c.SaveUploadedFile(file, path)
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

func SohoDownload(c *gin.Context) string {
	account, _ := c.Get("account")
	filename := c.Param("filename")
	param := c.Param("param")

	lowerFilename := strings.ToLower(filename)

	if strings.HasSuffix(lowerFilename, ".doc") || strings.HasSuffix(lowerFilename, ".pdf") || strings.HasSuffix(lowerFilename, ".ppt") || strings.HasSuffix(lowerFilename, ".jpf") ||
		strings.HasSuffix(lowerFilename, ".png") || strings.HasSuffix(lowerFilename, ".txt") || strings.HasSuffix(lowerFilename, ".gif") {

		var path string
		if param == "avatar" {
			path = os.Getenv("SOHO_WORK_PATH") + "/" + fmt.Sprintf("%v", account) + "/avatar/" + filename
		} else if param == "work" {
			path = os.Getenv("SOHO_WORK_PATH") + "/" + fmt.Sprintf("%v", account) + "/" + filename
		}

		return path
	} else {
		return ""
	}
}

func GetSohoWork(account string, param string) ([]model.SohoWork, error) {
	work := []model.SohoWork{}

	// 頭像
	where := "account = ? "
	if param == "avatar" {
		where += "and file_type ='1'"
	} else if param == "work" {
		where += "and file_type ='0'"
	}

	err := dao.GormSession.Select("*").Where(where, account).Find(&work).Error

	if err != nil {
		return nil, err
	} else {
		return work, nil
	}
}

func DeleteSohoWork(c *gin.Context) error {
	account, _ := c.Get("account")
	filename := c.Param("filename")
	param := c.Param("param")

	tx := dao.GormSession.Begin()

	work := model.SohoWork{
		Account:  fmt.Sprintf("%v", account),
		FileName: filename,
	}

	// 頭像
	if param == "avatar" {
		work.FileType = "1"
	} else if param == "work" {
		work.FileType = "0"
	}

	if err := tx.Delete(&work).Error; err != nil {
		tx.Rollback()
		return err
	}

	var path string
	if param == "avatar" {
		path = os.Getenv("SOHO_WORK_PATH") + "/" + fmt.Sprintf("%v", account) + "/avatar/" + filename
	} else {
		path = os.Getenv("SOHO_WORK_PATH") + "/" + fmt.Sprintf("%v", account) + "/" + filename
	}

	if err := os.Remove(path); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil

}
