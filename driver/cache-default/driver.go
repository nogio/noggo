package cache_default


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
)

type (
	DefaultCacheDriver struct {
	}
)

//返回驱动
func Driver() (driver.CacheDriver) {
	return &DefaultCacheDriver{}
}






//驱动连接
func (drv *DefaultCacheDriver) Connect(config Map) (driver.CacheConnect,error) {
	return &DefaultCacheConnect{
		config: config, caches: map[string]Any{},
	},nil
}
