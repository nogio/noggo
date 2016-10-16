package consts

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
)

func init() {

	//注册一些正则表达式
	noggo.Const.Regular(Map{

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
}
