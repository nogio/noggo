package https

import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
)

func init() {


	noggo.Http.Route("test", Map{
		"uri": "/",
		"route": Map{
			"name": "测试HTTP", "text": "测试HTTP",
			"action": func(ctx *noggo.HttpContext) {
				ctx.Text("index")
			},
		},
	})

}
