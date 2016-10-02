package middlers

import (
	"github.com/nogio/noggo"
	//"github.com/nogio/noggo/middler/plan-logger"
)

func init() {

	noggo.Plan.RequestFilter("asdfasdf", func(ctx *noggo.PlanContext) {
		ctx.Next()
	})


}
