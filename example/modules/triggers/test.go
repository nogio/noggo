package triggers

import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"fmt"
)

func init() {

	noggo.Trigger.Route("test", Map{
		"uri": "/test/teasdf",
		"name": "测试触发器", "text": "测试触发器说明",
		"action": func(ctx *noggo.TriggerCtx) {
			fmt.Println("测试触发器工作了")
			ctx.Finish();
		},
	})

}
