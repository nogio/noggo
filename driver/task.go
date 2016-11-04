/*
	接口
	2016-10-22  定稿
*/


package driver

import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/driver/task-default"
)



//默认任务引擎
func TaskDefault() (noggo.TaskDriver) {
	return task_default.Driver()
}
