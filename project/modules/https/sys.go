package https

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
)

func init() {


	noggo.Http.Route("_sys_.doc", Map {
		"uri": "/_sys_/doc",
		"route": Map {
			"name":	"自动生成系统文档", "text": "自动生成系统文档", "coded": true,
			"action": func(ctx *noggo.HttpContext) {

				ctx.Data["cryptos"] = noggo.Mapping.Cryptos()
				ctx.Data["states"] = noggo.Const.StateStrings()
				ctx.Data["models"] = noggo.Data.Models()
				ctx.Data["routes"] = ctx.Node.Http.Routes()

				ctx.View("_sys_/doc")
			},
		},
	})

}