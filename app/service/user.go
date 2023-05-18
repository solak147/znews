package service

import (
	"errors"
	"fmt"
	"znews/app/dao"
	"znews/app/middleware"
	"znews/app/model"

	"github.com/sirupsen/logrus"
)

var UserFields = []string{"account", "email"}

func GetUser(account string) (*model.User, error) {
	user := model.User{}
	err := dao.GormSession.Select("*").Where("account=?", account).First(&user).Error
	if err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

func GetUserByPwd(account string, password string) (*model.User, error) {
	user := model.User{}
	err := dao.GormSession.Select(UserFields).Where("account=? and password=?", account, password).First(&user).Error
	if err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

func UpdateUser(form model.ProfileSave) error {

	user := model.User{
		Name:         form.Name,
		Zipcode:      form.Zipcode,
		Phone:        form.Phone,
		Introduction: form.Introduction,
	}

	//是否勾選要修改密碼
	if form.PwdSwitch {
		u, e := GetUserByPwd(form.Account, form.OldPassword)
		if e != nil || u == nil {
			return errors.New("舊密碼錯誤")
		}
		user.Password = form.Password
	}

	err := dao.GormSession.Model(&model.User{}).Where("account = ?", form.Account).Updates(user).Error
	if err != nil {
		return err
	} else {
		return nil
	}

}

func Register(form model.RegisterStep3) error {
	if err := regexpRigister(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, form.Account); err != nil {
		return err
	}

	if err := regexpRigister(`^(?=.*[A-Z])(?=.*[a-z])(?=.*\d)\S{8,16}$`, form.Password); err != nil {
		return err
	}

	if err := regexpRigister(`^\S.{0,13}\S?$`, form.Name); err != nil {
		return err
	}

	if err := regexpRigister(`^\d{1,15}$`, form.Phone); err != nil {
		return err
	}

	if !CheckUserExit(form.Account) {
		return fmt.Errorf("User exists.")
	}

	user := model.User{
		Account:      form.Account,
		Password:     form.Password,
		Name:         form.Name,
		Zipcode:      form.Zipcode,
		Phone:        form.Phone,
		Email:        form.Account, //預設信箱同帳號
		Introduction: form.Introduction,
	}

	insertErr := dao.GormSession.Model(&model.User{}).Create(&user).Error
	return insertErr
}

func CheckUserExit(account string) bool {
	result := false
	var user model.User

	dbResult := dao.GormSession.Where("account = ?", account).Find(&user)
	if dbResult.Error != nil {
		middleware.Logger().WithFields(logrus.Fields{
			"name": "Get User Info Failed:",
		}).Error(dbResult.Error)
	} else {
		result = true
	}
	return result
}

func SohoSetting(account string, form model.SohoSettingForm) error {

	if err := regexpRigister(`^\S.{0,13}\S?$`, form.Name); err != nil {
		return err
	}

	if err := regexpRigister(`^\d{1,15}$`, form.Phone); err != nil {
		return err
	}

	if err := regexpRigister(`^\d{0,4}$`, form.CityTalk); err != nil {
		return err
	}

	if err := regexpRigister(`^\d{0,10}$`, form.CityTalk2); err != nil {
		return err
	}

	if err := regexpRigister(`^\d{0,5}$`, form.Extension); err != nil {
		return err
	}

	if err := regexpRigister(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, form.Email); err != nil {
		return err
	}

	if err := regexpRigister(`^[\s\S]{1,200}$`, form.Description); err != nil {
		return err
	}

	tx := dao.GormSession.Begin()

	sohoD := model.SohoSetting{}
	if err := tx.Where("account = ?", account).Delete(&sohoD).Error; err != nil {
		tx.Rollback()
		return err
	}

	set := model.SohoSetting{
		Account:     account,
		Open:        form.Open,
		Name:        form.Name,
		Role:        form.Role,
		Phone:       form.Phone,
		CityTalk:    form.CityTalk,
		CityTalk2:   form.CityTalk2,
		Extension:   form.Extension,
		Email:       form.Email,
		Zipcode:     form.Zipcode,
		Type:        form.Type,
		Exp:         form.ExpVal,
		Description: form.Description,
	}

	if err := tx.Model(&model.SohoSetting{}).Create(&set).Error; err != nil {
		tx.Rollback()
		return err
	} else {
		tx.Commit()
		return nil
	}
}

func ChkSohoSetting(account string) error {
	var cnt int64
	if err := dao.GormSession.Model(&model.SohoSetting{}).Where("account = ?", account).Count(&cnt).Error; err != nil {
		return err
	} else {
		if cnt == 0 {
			return errors.New("尚未填寫接案設定")
		}
		return nil
	}
}

func SohoSettingInit(account string) (interface{}, error) {

	set := model.SohoSetting{}
	if err := dao.GormSession.Select("*").Where("account=?", account).First(&set).Error; err != nil {

		if err.Error() == "record not found" {

			user := model.User{}
			if err := dao.GormSession.Select("name,phone,zipcode,email").Where("account=?", account).First(&user).Error; err != nil {
				return nil, err
			}
			return user, nil

		} else {
			return nil, err
		}

	}
	return set, nil
}

// 新增作品網址
func AddSohoUrl(account string, form model.SohoUrlForm) error {

	soho := model.SohoUrl{
		Account: account,
		Url:     form.Url,
		Explain: form.Explain,
	}

	if err := dao.GormSession.Model(&model.SohoUrl{}).Create(&soho).Error; err != nil {
		return err
	} else {
		return nil
	}

}

// 取得作品網址
func GetSohoUrl(account string) ([]model.SohoUrl, error) {

	soho := []model.SohoUrl{}
	if err := dao.GormSession.Select("*").Where("account=?", account).Find(&soho).Error; err != nil {
		return nil, err
	} else {
		return soho, nil
	}

}

// 刪除作品網址
func DeleteSohoUrl(account string, url string) error {

	soho := model.SohoUrl{}
	if err := dao.GormSession.Where("account = ? AND url = ?", account, url).Delete(&soho).Error; err != nil {
		return err
	}
	return nil

}
