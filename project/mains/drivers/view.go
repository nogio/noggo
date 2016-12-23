package drivers

import (
    "github.com/nogio/noggo"
    "github.com/nogio/noggo/driver/view-default"
)

func init() {
    //默认驱动
    noggo.Driver("default", view_default.Driver())
}
