/*
	缓存接口
	2016-10-22  未定稿
	因为返回的时候，其实无法确定是什么任务，也无法强制转换
	所有只能在使用缓存的地方，手动设置和转换类型
	所以暂时不能确认接口文档
*/


package driver

import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/driver/cache-default"
	"github.com/nogio/noggo/driver/cache-redis"
)



//默认缓存引擎
func CacheDefault() (noggo.CacheDriver) {
	return cache_default.Driver()
}


//redis缓存引擎
func CacheRedis() (noggo.CacheDriver) {
	return cache_redis.Driver()
}
