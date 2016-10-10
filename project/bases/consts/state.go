package consts

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
)

func init() {

	//注册一些状态
	noggo.Const.State(Map{

		"ok": 0,
		"no": 1,

		"none": 2,  //拒绝访问

		"found":	3,  //不存在
		"error":	4,  //有错误
		"failed":	5, 	//失败的
		"denied":	6,	//拒绝的

		"coding": 7,	//编码中

		"map.empty": 8,  //参数不能为空
		"map.error": 9,  //参数类型错误

		"auth.empty": 10,  //auth签名？不存在？
		"auth.error": 11,  //auth数据？不存在？

		"item.empty": 12,    //item参数不存在
		"item.error": 13,    //item数据不存在

		"data.empty": 14,    //返回数据 不可为空
		"data.error": 15,    //返回数据 生成失败

	})
}
