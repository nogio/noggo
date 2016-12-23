package utils

import (
	"os"
	"io"
)


func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//文件复制
func FileCopy(src, des string) bool {
	srcFile, err := os.Open(src)
	if err == nil {
		defer srcFile.Close()

		//创建新文件
		desFile, err := os.Create(des)
		if err == nil {
			defer desFile.Close()

			//开始复制
			_,err := io.Copy(desFile, srcFile)
			if err == nil {
				return true
			}
		}
	}
	return false
}