package noggo

import (
	. "github.com/nogio/noggo/base"
	"sync"
	"errors"
)



// 缓存驱动接口定义 begin
type (
	//缓存驱动
	CacheDriver interface {
		Connect(config Map) (CacheConnect,error)
	}
	//缓存连接
	CacheConnect interface {
		//打开连接
		Open() error
		//关闭连接
		Close() error
		//获取数据库对象
		Base(string) (CacheBase,error)
	}

	//缓存库
	CacheBase interface {
		Close() error


		//查询缓存，自带值包装函数
		Get(key string) (Any,error)
		//更新缓存数据，不存在则创建，存在就更新
		Set(key string, val Any, exp int64) error
		//删除缓存
		Del(key string) error
		//清空缓存，如传prefix则表示清空固定前缀的缓存
		Empty(prefixs ...string) (error)
		//获取keys
		Keys(prefixs ...string) ([]string,error)

	}
)
// 缓存驱动接口定义 end









type (
	cacheGlobal struct {
		mutex       sync.Mutex
		drivers     map[string]CacheDriver

		connects    map[string]CacheConnect
	}
)



//注册缓存驱动
func (global *cacheGlobal) Driver(name string, config CacheDriver) {
	global.mutex.Lock()
	defer global.mutex.Unlock()


	if global.drivers == nil {
		global.drivers = map[string]CacheDriver{}
	}

	if config == nil {
		panic("缓存: 驱动不可为空")
	}
	global.drivers[name] = config
}


//连接驱动
func (global *cacheGlobal) connect(config *cacheConfig) (CacheConnect,error) {
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
func (global *cacheGlobal) Base(name string) (CacheBase) {
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
func (base *noCacheBase) Close() (error) {
	return nil
}
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
