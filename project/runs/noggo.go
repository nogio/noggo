package main

import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	_ "../bases"
	_ "../drivers"
	_ "../globals"
	_ "../nodes"
	"time"
)

func main() {
	noggo.Init()

	noggo.Logger.Debug("哈哈哈哈哈")
	noggo.Trigger.Touch("test", Map{ "id": 1 })

	noggo.Task.Touch("test", time.Second * 10, Map{ "id": 10 })
	noggo.Task.Touch("test", time.Second * 30, Map{ "id": 10 })

	nog := noggo.New("www")
	nog.Run()

	noggo.Exit()
}
