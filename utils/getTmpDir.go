package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// 获取临时目录
func GetTmpDir() string {
	_dir, _ := filepath.Abs("./tmp")
	exist, err := PathExists(_dir)
	if err != nil {
		fmt.Printf("Get tmp dir failed [%v]", err)
		return ""
	}

	if !exist {
		err = os.Mkdir(_dir, os.ModePerm)
		if err != nil {
			fmt.Printf("Create tmp dir failed [%v]", err)
			return ""
		}
	}

	return _dir
}

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
