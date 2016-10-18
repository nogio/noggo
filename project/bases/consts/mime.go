package consts

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
)

func init() {

	//注册所有的mime类型
	noggo.Const.Mime(Map{
		"text": "text/explain",
		"html": "text/html",
		"xml": "text/xml",
		"json": "text/json",
		"file": "text/explain",
		"down": "text/explain",
		"script": "application/x-javascript",
		"view": "text/html",

		"jpg": "image/jpeg",
		"gif": "image/gif",

		"test": "type/test",
	})
}
