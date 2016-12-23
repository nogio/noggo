package drivers

import (
    "github.com/nogio/noggo"
    "github.com/nogio/noggo/driver/plan-default"
)

func init() {
    //默认驱动
    noggo.Driver("default", plan_default.Driver())
}
