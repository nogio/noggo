package drivers


import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/driver/task-default"
)


func init() {
	//注册默认任务驱动
	noggo.Task.Driver("default", task_default.Driver())
}
