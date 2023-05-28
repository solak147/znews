package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"znews/app/model"

	"github.com/dlclark/regexp2"
)

func regexpRigister(pattern string, matchStr string) error {
	// 编译正则表达式
	re := regexp2.MustCompile(pattern, 0)

	isMatch, err := re.MatchString(matchStr)
	if err == nil && isMatch {
		return nil
	} else {
		errStr := matchStr + " not match regexp"

		if err != nil {
			errStr = errStr + ", err : " + err.Error()
		}
		return errors.New(errStr)
	}

}

func cryptoSha256(password string) string {
	myCrypto := os.Getenv("SHA256")
	password = password + myCrypto

	// 创建哈希对象
	hash := sha256.New()

	// 将用户输入的密码转换为字节数组并计算哈希值
	hash.Write([]byte(password))
	userInputPasswordHash := hash.Sum(nil)

	// 将用户输入密码的哈希值转换为十六进制字符串
	userInputPasswordHashStr := hex.EncodeToString(userInputPasswordHash)

	return userInputPasswordHashStr
}

func chkCaseForm(form model.CaseForm) error {
	if err := regexpRigister(`^.{3,20}$`, form.Title); err != nil {
		return err
	}

	if err := regexpRigister(`^.{4}$`, form.Type); err != nil {
		return err
	}

	if err := regexpRigister(`^[o,i]$`, form.Kind); err != nil {
		return err
	}

	// 0-30 or yyyy/mm/dd
	if err := regexpRigister(`^([0-9]{4}/(0[1-9]|1[012])/(0[1-9]|[12][0-9]|3[01])|([12][0-9]|30|[0-9])|)$`, form.ExpectDate); err != nil {
		return err
	}

	if err := regexpRigister(`^[1,2,3]$`, form.ExpectDateChk); err != nil {
		return err
	}

	if err := regexpRigister(`^[\s\S]{1,200}$`, form.WorkContent); err != nil {
		return err
	}

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

	if err := regexpRigister(`^[a-zA-Z0-9_-]*$`, form.Line); err != nil {
		return err
	}

	return nil
}
