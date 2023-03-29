package service

import (
	"fmt"
	"znews/app/dao"
	"znews/app/model"
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

func RegisterOneUser(account string, password string, email string) error {
	if !CheckOneUser(account) {
		return fmt.Errorf("User exists.")
	}
	user := model.User{
		Account:  account,
		Password: password,
		Email:    email,
	}
	insertErr := dao.SqlSession.Model(&model.User{}).Create(&user).Error
	return insertErr
}

func CheckOneUser(account string) bool {
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
