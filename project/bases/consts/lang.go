package consts

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
)

func init() {

	//注册默认的语言配置
	noggo.Const.Lang("default", Map{

		"ok":	"成功",
		"no":	"失败",

		"test":	"测试",

	})
}
