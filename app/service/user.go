package service

import (
	"errors"
	"fmt"
	"znews/app/dao"
	"znews/app/model"

	"github.com/dlclark/regexp2"
)

var UserFields = []string{"id", "account", "email"}

func SelectOneUsers(id int64) (*model.User, error) {
	userOne := &model.User{}
	err := dao.SqlSession.Select(UserFields).Where("id=?", id).First(&userOne).Error
	if err != nil {
		return nil, err
	} else {
		return userOne, nil
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
		fmt.Printf("Get User Info Failed:%v\n", dbResult.Error)
	} else {
		result = true
	}
	return result
}

func GetUser(account string, password string) (*model.User, error) {
	userOne := &model.User{}
	err := dao.SqlSession.Select(UserFields).Where("account=? and password=?", account, password).First(&userOne).Error
	if err != nil {
		return nil, err
	} else {
		return userOne, nil
	}
}

func regexpRigister(pattern string, matchStr string) error {
	// 编译正则表达式
	re := regexp2.MustCompile(pattern, 0)

	isMatch, err := re.MatchString(matchStr)
	if err == nil && isMatch {
		return nil
	} else {
		return errors.New(matchStr + "not match regexp, err:" + err.Error())
	}

}
