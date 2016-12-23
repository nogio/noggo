package drivers

import (
    "github.com/nogio/noggo"
    "github.com/nogio/noggo/driver/task-default"
)

func init() {
    //默认驱动
    noggo.Driver("default", task_default.Driver())
}
