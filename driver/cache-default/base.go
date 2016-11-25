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
//关闭库
func (base *DefaultCacheBase) Close() (error) {
	return nil
}

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
func (base *DefaultCacheBase) Set(key string, val Any, exps ...int64) (error) {
	base.conn.mutex.Lock()
	defer base.conn.mutex.Unlock()

	//暂时不考虑过期的问题
	base.conn.caches[key] = val

	return nil
}
//删除数据
func (base *DefaultCacheBase) Del(key string) (error) {
	base.conn.mutex.Lock()
	defer base.conn.mutex.Unlock()

	delete(base.conn.caches, key)

	return nil
}



//计数器
func (base *DefaultCacheBase) Num(key string, nums ...int) (int64,error) {
	base.conn.mutex.Lock()
	defer base.conn.mutex.Unlock()

	num := 1
	if len(nums) > 0 {
		num = nums[0]
	}

	if v,ok := base.conn.caches[key].(int64); ok {
		val := v+int64(num)
		base.conn.caches[key] = val

		return val, nil
	} else {
		val := int64(num)
		base.conn.caches[key] = val

		return val,nil
	}
}








//清空数据
func (base *DefaultCacheBase) Clear(prefixs ...string) (error) {
	base.conn.mutex.Lock()
	defer base.conn.mutex.Unlock()

	keys,err := base.Keys(prefixs...)
	if err == nil {
		for _,key := range keys {
			delete(base.conn.caches, key)
		}
	}
	return nil
}



//获取keys
//暂时不支持前缀查询
func (base *DefaultCacheBase) Keys(prefixs ...string) ([]string,error) {
	base.conn.mutex.RLock()
	defer base.conn.mutex.RUnlock()

	keys := []string{}
	for k,_ := range base.conn.caches {
		keys = append(keys, k)
	}

	return keys,nil
}
