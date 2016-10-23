package cache_default


import (
	. "github.com/nogio/noggo/base"
	"errors"
)

type (
	DefaultCacheBase struct {
		name    string
		conn    *DefaultCacheConnect
	}
)

//获取数据
//数据应该复制一份，要不然写的时候会有同步问题
//因为Any相当于引用
func (base *DefaultCacheBase) Get(key string) (Any,error) {
	base.conn.mutex.RLock()
	defer base.conn.mutex.RUnlock()

	if v,ok := base.conn.caches[key]; ok {
		return v, nil
	}

	return nil,errors.New("无缓存")
}
//设置数据
func (base *DefaultCacheBase) Set(key string, val Any, expiry int64) (error) {
	base.conn.mutex.Lock()
	defer base.conn.mutex.Unlock()

	//暂时不直接过期的问题
	base.conn.caches[key] = val

	return nil
}
//删除数据
func (base *DefaultCacheBase) Del(key string) (error) {
	base.conn.mutex.Lock()
	defer base.conn.mutex.Unlock()

	//暂时不直接过期的问题
	delete(base.conn.caches, key)

	return nil
}
//空清数据
func (base *DefaultCacheBase) Empty(prefixs ...string) (error) {

	keys,err := base.Keys(prefixs...)
	if err == nil {
		base.conn.mutex.Lock()
		for _,key := range keys {
			delete(base.conn.caches, key)
		}
		base.conn.mutex.Unlock()
	}
	return nil
}

//获取keys
func (base *DefaultCacheBase) Keys(prefixs ...string) ([]string,error) {
	base.conn.mutex.RLock()
	defer base.conn.mutex.RUnlock()

	keys := []string{}
	for k,_ := range base.conn.caches {
		keys = append(keys, k)
	}

	return keys,nil
}
