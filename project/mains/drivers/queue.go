package drivers

import (
    "github.com/nogio/noggo"
    "github.com/nogio/noggo/driver/queue-default"
)

func init() {
    //默认驱动
    noggo.Driver("default", queue_default.Driver())
}
