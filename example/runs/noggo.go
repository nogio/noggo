package main

import (
	"fmt"
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	_ "../drivers"
	_ "../modules/triggers"
)

func main() {

	fmt.Printf("config=%v\n", noggo.Config)


	//触发test
	noggo.Trigger.Touch("test", Map{ "msg": 111 })
	//noggo.Trigger.Touch("test", Map{ "msg": 222 })
	//noggo.Trigger.Touch("test", Map{ "msg": 333 })

	//因为三个请求是同步的， session是内存的，读的是同一个session
	//同时写。 就出错了。。  session-memory驱动要处理一下这样同步的问题
	//返回来的， 应该是一个副本， 这样本地随便写，  Update的时候才写回去

	noggo.Launch()
}
