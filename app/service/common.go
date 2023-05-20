package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"

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
