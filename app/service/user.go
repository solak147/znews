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

	insertErr := dao.GormSession.Model(&model.SohoSetting{}).Create(&set).Error
	if insertErr != nil {
		return insertErr
	} else {
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
