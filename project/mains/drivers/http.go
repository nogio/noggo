package drivers

import (
    "github.com/nogio/noggo"
    "github.com/nogio/noggo/driver/http-default"
)

func init() {
    //默认驱动
    noggo.Driver("default", http_default.Driver())
}
