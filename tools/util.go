package tools

import (
	"errors"
	"os"
	"regexp"
	"unicode"
)

func PathFileExists(p string, ignoreExists bool) error {
	exists, _ := FileExists(p)
	if exists && ignoreExists == false {
		err := errors.New("folder exists")
		return err
	}
	if !exists {
		err := os.MkdirAll(p, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// IsChineseChar 判断是否是中文
func IsChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) || (regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(r))) {
			return true
		}
	}
	return false
}

func JudgeType2() {

}
