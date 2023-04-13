package service

import (
	"errors"

	"github.com/dlclark/regexp2"
)

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
