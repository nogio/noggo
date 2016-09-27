package consts

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
)

func init() {

	noggo.Const.Mime(Map{
		"text": "text/explain",
		"html": "text/html",
		"xml": "text/xml",
		"json": "application/json",
		"script": "application/x-javascript",

		"jpg": "image/jpeg",
		"gif": "image/gif",

		"test": "type/test",
	})
}
