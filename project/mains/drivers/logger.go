package drivers

import (
    "github.com/nogio/noggo"
    "github.com/nogio/noggo/driver/logger-default"
)

func init() {
    //默认驱动
    noggo.Driver("default", logger_default.Driver())
}
