package util

// author：xiaosheng

import (
	"errors"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"
	"xiaosheng/tools"
)

// PathFileExists 目录不存在则创建目录
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

// FileExists 判断文件是否存在
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

// FileCreate 文件不存在则创建
func FileCreate(p string, flag bool) string {
	exists, _ := FileExists(p)
	PathFileExists(tools.ConfigPath, false)
	if exists {
		// 有这个文件就添加后缀
		if flag {
			split := strings.Split(p, ".")
			allWorld := ""
			length := len(split)
			for index, item := range split {
				if index == length-2 {
					allWorld = allWorld + item + "_new"
				} else if index == length-1 {
					allWorld = allWorld + "." + item
				} else {
					allWorld = allWorld + item + "."
				}
			}
			p = allWorld
			f, err := os.Create(p)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			return p
		}
		return ""

	} else {
		f, err := os.Create(p)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	}
	return ""
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

//func DealWithName() {
//
//}
