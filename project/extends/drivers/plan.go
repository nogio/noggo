package drivers


import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/driver/plan-default"
)


func init() {
	//注册默认计划驱动
	noggo.Plan.Driver("default", plan_default.Driver())
}
