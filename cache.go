package noggo

import (
	. "github.com/nogio/noggo/base"
	"sync"
	"errors"
	"time"
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

		//读缓存，自带值包装函数
		Get(key string) (Any,error)
		//写缓存数据，不存在则创建，存在就更新
		Set(key string, val Any, exps ...int64) error
		//删除缓存
		Del(key string) error

		//计数incr
		Num(key string,nums ...int) (int64,error)

		//获取keys
		Keys(prefixs ...string) ([]string,error)

		//清理缓存，如传prefix则表示清空固定前缀的缓存
		Clear(prefixs ...string) (error)

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






//缓存值打包
func (global *cacheGlobal) Pack(val Any) (Map,error) {

	valType := "string"

	//获取值类型
	switch val.(type) {
	case string:
		valType = "string"
	case []string:
		valType = "[string]"

	case int,int8,int16,int32,int64:
		valType = "int"
	case []int,[]int8,[]int16,[]int32,[]int64:
		valType = "[int]"

	case float32,float64:
		valType = "float"
	case []float32,[]float64:
		valType = "[float]"

	case bool:
		valType = "bool"
	case []bool:
		valType = "[bool]"

	case time.Time:
		valType = "time"
	case []time.Time:
		valType = "[time]"

	case Map,map[string]interface{}:
		valType = "json"
	case []Map,[]map[string]interface{}:
		valType = "[json]"
	}

	return Map{
		"type": valType, "data": val,
	},nil
}
//缓存值解包
func (global *cacheGlobal) UnPack(val Map) (Any,error) {

	if val != nil && val["type"] != nil && val["data"] != nil{

		if tt,ok := val["type"].(string); ok {

			config := Map{
				"data": Map{ "type": tt, "must": true },
			}
			value := Map{}

			err := Mapping.Parse([]string{}, config, val, value, false, false)
			if err != nil {
				return nil,err
			} else {
				return value["data"],nil
			}
		}
	}

	return nil,errors.New("解包失败")
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
func (base *noCacheBase) Set(key string, val Any, exps ...int64) (error) {
	return errors.New("无缓存")
}
func (base *noCacheBase) Del(key string) (error) {
	return errors.New("无缓存")
}
func (base *noCacheBase) Num(key string, nums ...int) (int64,error) {
	return int64(0),errors.New("无缓存")
}
func (base *noCacheBase) Clear(prefixs ...string) (error) {
	return errors.New("无缓存")
}
func (base *noCacheBase) Keys(prefixs ...string) ([]string,error) {
	return []string{},errors.New("无缓存")
}
