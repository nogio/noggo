package cache_default


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
)

type (
	DefaultCacheDriver struct {
	}
)

//返回驱动
func Driver() (noggo.CacheDriver) {
	return &DefaultCacheDriver{}
}






//驱动连接
func (drv *DefaultCacheDriver) Connect(config Map) (noggo.CacheConnect,error) {
	return &DefaultCacheConnect{
		config: config, caches: map[string]Any{},
	},nil
}
