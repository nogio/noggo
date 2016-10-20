package native

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
)

func init() {

	//注册mime类型
	noggo.Const.Mimes(Map{

		"text": "text/explain",
		"html": "text/html",
		"xml": "text/xml",
		"json": "application/json",
		"file": "application/octet-stream",
		"down": "application/octet-stream",
		"script": "application/x-javascript",
		"view": "text/html",

		"jpg": "image/jpeg",
		"gif": "image/gif",
		"test": "type/test",
	})


	//注册状态
	noggo.Const.States(Map{
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


	//注册正则

	noggo.Const.Regulars(Map{

		"password": `^[0-9A-Fa-f]{32}$`,

		"number": `^[0-9]+$`,
		"float": `^[0-9]+$`,

		"date": []string{
			`^(\d{4})-(\d{2})-(\d{2})$`,
			`^(\d{10, 15})$`,
		},
		"datetime": []string{
			`^(\d{4})-(\d{2})-(\d{2})$`,
			`^(\d{4})-(\d{2})-(\d{2}) (\d{2}):(\d{2}):(\d{2})$`,
			`^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2}):(\d{2})Z$`,
			`^(\d{4})-(\d{2})-(\d{2}) (\d{2}):(\d{2}):(\d{2})\.(\d{3})$`,
			`^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2}):(\d{2})\.(\d{3})Z$`,
			`^(\d{10, 15})$`,
		},

		"mobile": `^1[0-9]{10}$`,
		"email": `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`,
		"array": []string{
			"^asdfasfasf/",
		},
		"test": `^test/`,
	})



	//注册多语言字串
	noggo.Const.Langs(noggo.ConstLangDefault, Map{
		"ok": "成功",
		"no": "失败",

		"none": "拒绝访问",

		"found":  "前方有坑",
		"error":  "系统错误",
		"failed": "请求失败",
		"denied": "请求拒绝",

		"coding": "编码中",

		"map.empty": "%s不可为空",
		"map.error": "%s不是有效的值",

		"auth.empty": "%s未登录",
		"auth.error": "%s未登录",

		"item.empty": "%s参数不存在",
		"item.error": "%s记录不存在",

		"data.empty": "%s数据不可为空",
		"data.error": "%s数据生成失败",

	})


}
