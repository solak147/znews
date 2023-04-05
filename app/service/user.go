package service

import (
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
	regexpRigister(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, form.Account)
	regexpRigister(`^(?=.*[A-Z])(?=.*[a-z])(?=.*\d)\S{8,16}$`, form.Password)
	regexpRigister(`^\S.{0,13}\S?$`, form.Name)
	regexpRigister(`^\d{1,15}$`, form.Phone)

	if !CheckUserExit(form.Account) {
		return fmt.Errorf("User exists.")
	}

	form.Email = form.Account //預設信箱同帳號

	insertErr := dao.SqlSession.Model(&model.User{}).Create(&form).Error
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

func regexpRigister(pattern string, matchStr string) bool {
	// 编译正则表达式
	re := regexp2.MustCompile(pattern, 0)

	isMatch, err := re.MatchString(matchStr)
	if err != nil && isMatch {
		return true
	} else {
		fmt.Println(matchStr + "not match regexp, err:" + err.Error())
		return false
	}

}
