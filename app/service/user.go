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
	user := &model.User{}
	err := dao.SqlSession.Select("*").Where("account=?", account).First(&user).Error
	if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func GetUserByPwd(account string, password string) (*model.User, error) {
	user := &model.User{}
	err := dao.SqlSession.Select(UserFields).Where("account=? and password=?", account, password).First(&user).Error
	if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func UpdateUser(form model.ProfileSave) error {

	user := model.User{
		Name:         form.Name,
		Zipcode:      form.Zipcode,
		Phone:        form.Phone,
		Introduction: form.Introduction,
	}

	if form.PwdSwitch {
		u, e := GetUserByPwd(form.Account, form.OldPassword)
		if e != nil || u == nil {
			return errors.New("舊密碼錯誤")
		}
		user.Password = form.Password
	}

	err := dao.SqlSession.Model(&model.User{}).Where("account = ?", form.Account).Updates(user).Error
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

	insertErr := dao.SqlSession.Model(&model.User{}).Create(&user).Error
	return insertErr
}

func CheckUserExit(account string) bool {
	result := false
	var user model.User

	dbResult := dao.SqlSession.Where("account = ?", account).Find(&user)
	if dbResult.Error != nil {
		middleware.Logger().WithFields(logrus.Fields{
			"name": "Get User Info Failed:",
		}).Error(dbResult.Error)
	} else {
		result = true
	}
	return result
}
