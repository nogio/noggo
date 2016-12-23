package drivers

import (
    "github.com/nogio/noggo"
    "github.com/nogio/noggo/driver/session-default"
    "github.com/nogio/noggo/driver/session-redis"
)

func init() {
    //默认驱动
    noggo.Driver("default", session_default.Driver())
    noggo.Driver("redis", session_redis.Driver())
}
