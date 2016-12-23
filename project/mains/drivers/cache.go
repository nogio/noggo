package drivers

import (
    "github.com/nogio/noggo"
    "github.com/nogio/noggo/driver/cache-default"
    "github.com/nogio/noggo/driver/cache-redis"
)

func init() {
    //默认驱动
    noggo.Driver("default", cache_default.Driver())
    noggo.Driver("redis", cache_redis.Driver())
}
