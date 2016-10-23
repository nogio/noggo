package noggo

import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"sync"
	"errors"
)

type (
	cacheGlobal struct {
		mutex       sync.Mutex
		drivers     map[string]driver.CacheDriver

		connects    map[string]driver.CacheConnect
	}
)



//注册缓存驱动
func (global *cacheGlobal) Driver(name string, config driver.CacheDriver) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.drivers == nil {
		global.drivers = map[string]driver.CacheDriver{}
	}

	if config == nil {
		panic("缓存: 驱动不可为空")
	}
	global.drivers[name] = config
}


//连接驱动
func (global *cacheGlobal) connect(config *cacheConfig) (driver.CacheConnect,error) {
	if cacheDriver,ok := global.drivers[config.Driver]; ok {
		return cacheDriver.Connect(config.Config)
	} else {
		panic("缓存：不支持的驱动 " + config.Driver)
	}
}

//缓存初始化
func (global *cacheGlobal) init() {
	global.initCache()
}

//初始化缓存驱动
func (global *cacheGlobal) initCache() {
	for name,config := range Config.Cache {
		con,err := global.connect(config)
		if err != nil {
			panic("数据：连接失败：" + err.Error())
		} else {
			err := con.Open()
			if err != nil {
				panic("数据：打开连接失败：" + err.Error())
			} else {
				//保存连接
				global.connects[name] = con
			}
		}
	}
}
//缓存退出
func (global *cacheGlobal) exit() {
	global.exitCache()
}
//任务退出，缓存
func (global *cacheGlobal) exitCache() {
	for _,con := range global.connects {
		con.Close()
	}
}










//返回缓存Base对象
func (global *cacheGlobal) Base(name string) (driver.CacheBase) {
	if conn,ok := global.connects[name]; ok {
		db,err := conn.Base(name)
		if err == nil {
			return db
		}
	}
	return &noCacheBase{}
}




//---------------------------------------------------------------------------------

type (
	noCacheBase struct {}
)
func (base *noCacheBase) Get(key string) (Any,error) {
	return nil,errors.New("无缓存")
}
func (base *noCacheBase) Set(key string, val Any, expiry int64) (error) {
	return errors.New("无缓存")
}
func (base *noCacheBase) Del(key string) (error) {
	return errors.New("无缓存")
}
func (base *noCacheBase) Empty(prefixs ...string) (error) {
	return errors.New("无缓存")
}
func (base *noCacheBase) Keys(prefixs ...string) ([]string,error) {
	return []string{},errors.New("无缓存")
}
