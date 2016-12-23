package drivers

import (
    "github.com/nogio/noggo"
    "github.com/nogio/noggo/driver/event-default"
)

func init() {
    //默认驱动
    noggo.Driver("default", event_default.Driver())
}
