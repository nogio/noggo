/*
	计划接口
	2016-10-22  定稿
*/

package driver

import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/driver/plan-default"
)

//默认任务驱动
func PlanDefault() (noggo.PlanDriver) {
	return plan_default.Driver()
}
