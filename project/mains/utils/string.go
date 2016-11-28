package utils

import (
)



//密码加密格式
func Password(str string) string {
	return Md5(str)
}

