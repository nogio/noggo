package https

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
)

func init() {

	noggo.Http.Route("index", Map{
		"uri": "/",
		"route": Map{
			"name": "首页", "text": "首页",
			"action": func(ctx *noggo.HttpContext) {
				ctx.View("index")
			},
		},
	})


	noggo.Http.Route("test", Map{
		"uri": "/test",
		"route": Map{
			"name": "test", "text": "test",
			"action": func(ctx *noggo.HttpContext) {

				db := noggo.Data.Base("main"); defer db.Close()

				/*
				items,_ := db.Model("test").Limit(0, 20, Map{
					"id": Map{ ">" : 5000 },
				}, Map{
					"id": driver.DESC, "created": driver.ASC,
				})
				*/

				items,err := db.Model("test").Group("title")
				noggo.Logger.Debug("group", items, err)


				ctx.Text(items)
			},
		},
	})

}