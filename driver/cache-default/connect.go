package cache_default

import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"sync"
)

type (
	//数据库连接
	DefaultCacheConnect struct {
		config  Map
		mutex   sync.RWMutex
		caches  map[string]Any
	}
)

//打开连接
func (conn *DefaultCacheConnect) Open() error {
	return nil
}
//关闭连接
func (conn *DefaultCacheConnect) Close() error {
	return nil
}




func (conn *DefaultCacheConnect) Base(name string) (noggo.CacheBase,error) {
	return &DefaultCacheBase{name, conn},nil
}