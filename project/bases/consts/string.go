package consts

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
)

func init() {

	noggo.Const.Lang(Map{

		"ok":	"成功",
		"no":	"失败",

		"test":	"测试",

	}, "default")
}
