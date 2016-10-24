package noggo

import (
	. "github.com/nogio/noggo/base"
	"sync"
	"github.com/nogio/noggo/driver"
	"errors"
	"database/sql"
)

type (
	dataGlobal    struct {
		mutex       sync.Mutex
		drivers     map[string]driver.DataDriver

		//数据连接
		connects    map[string]driver.DataConnect

		models      map[string]map[string]Map
	}
)


//连接驱动
func (global *dataGlobal) connect(config *dataConfig) (driver.DataConnect,error) {
	if dataDriver,ok := global.drivers[config.Driver]; ok {
		return dataDriver.Connect(config.Config)
	} else {
		panic("数据：不支持的驱动 " + config.Driver)
	}
}


//注册数据驱动
func (global *dataGlobal) Driver(name string, config driver.DataDriver) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if config == nil {
		panic("数据: 驱动不可为空")
	}
	//不做存在判断，因为要支持后注册的驱动替换已注册的驱动
	//框架有可能自带几种默认驱动，并且是默认注册的，用户可以自行注册替换
	global.drivers[name] = config
}










//数据初始化
func (global *dataGlobal) init() {
	global.initData()
}
func (global *dataGlobal) initData() {

	//遍历数据配置
	for name,config := range Config.Data {
		conn,err := global.connect(config)
		if err != nil {
			panic("数据：连接失败：" + err.Error())
		} else {
			err := conn.Open()
			if err != nil {
				panic("数据：打开连接失败：" + err.Error())
			} else {

				//注册模型
				//先全局
				for k,v := range global.models[ConstNodeGlobal] {
					conn.Model(k, v)
				}
				//再库
				for k,v := range global.models[name] {
					conn.Model(k, v)
				}

				//保存连接
				global.connects[name] = conn

			}
		}
	}
}

//数据退出
func (global *dataGlobal) exit() {
	global.exitData()
}
func (global *dataGlobal) exitData() {
	for _,conn := range global.connects {
		conn.Close()
	}
}
























//注册模型
func (global *dataGlobal) Model(name string, config Map) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.models == nil {
		global.models = map[string]map[string]Map{}
	}

	//节点
	nodeName := ConstNodeGlobal
	if Current != "" {
		nodeName = Current
	}

	//如果节点配置不存在，创建
	if global.models[nodeName] == nil {
		global.models[nodeName] = map[string]Map{}
	}


	//可以后注册重写原有路由配置，所以直接保存
	global.models[nodeName][name] = config
}


//查询某库所有模型
func (global *dataGlobal) Models(bases ...string) (Map) {

	if len(bases) > 0 {
		base := bases[0]
		if global.models[base] != nil {

			m := Map{}
			for k,v := range global.models[base] {
				m[k] = v
			}
			return m
		}
	} else {
		m := Map{}
		for k,v := range global.models {
			m[k] = v
		}
		return m
	}

	return Map{}
}



//查询某表所有字段
func (global *dataGlobal) Fields(base, model string, maps ...Map) (Map) {

	m := Map{}

	if dc,ok := global.models[base]; ok {
		if mc,ok := dc[model]; ok {
			if fc,ok := mc["fields"]; ok {
				//不可以直接给,要给一个新的,要不么返回了引用,改了后, 原定义也改了
				for k,v := range fc.(Map) {
					m[k] = v
				}
			}
		}
	}

	//覆盖的map
	if len(maps) > 0 {
		for k,v := range maps[0] {
			m[k] = v
		}
	}

	return m
}

//查询某表部分字段
func (global *dataGlobal) Field(base, model string, fields []string, maps ...Map) Map {

	m := Map{}

	if dc,ok := global.models[base]; ok {
		if mc,ok := dc[model]; ok {
			if fc,ok := mc["fields"]; ok {

				// 后续考虑支持多级
				// fields中名称是user.avatar.id 这样的  当是mongodb时，就比较重要了
				for _,n := range fields {

					//字段是否存在
					if v,ok := fc.(Map)[n]; ok {
						m[n] = v
					}

				}

			}
		}
	}

	//覆盖的map
	if len(maps) > 0 {
		for k,v := range maps[0] {
			m[k] = v
		}
	}

	return m
}





































//返回DB对象
func (global *dataGlobal) Base(name string) (driver.DataBase) {
	if conn,ok := global.connects[name]; ok {

		//缓存相关
		var cb driver.CacheBase = nil
		if cfg,ok := Config.Data[name]; ok {
			if cfg.Cache != "" {
				cb = Cache.Base(cfg.Cache)
			}
		}

		db,err := conn.Base(name, cb)
		if err == nil {
			//返回
			return db
		}
	}

	return &noDataBase{}
}





//----------------------------------------------------------------------

type (
	noDataBase struct {}
	noDataModel struct {}
)
func (base *noDataBase) Close() {
}
func (base *noDataBase) Model(name string) (driver.DataModel) {
	return &noDataModel{}
}
func (base *noDataBase) Begin() (driver.DataBase) {
	return base
}
func (base *noDataBase) Submit() (error) {
	return errors.New("无数据")
}
func (base *noDataBase) Cancel() (error) {
	return errors.New("无数据")
}


//----------------------------------------------------





//Exec
func (base *noDataBase) Exec(query string, args ...interface{}) (sql.Result,error) {
	return nil,errors.New("无数据")
}
//Prepare
func (base *noDataBase) Prepare(query string) (*sql.Stmt, error) {
	return nil,errors.New("无数据")
}
//Query
func (base *noDataBase) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil,errors.New("无数据")
}
//QueryRow
func (base *noDataBase) QueryRow(query string, args ...interface{}) (*sql.Row) {
	return nil
}
//QueryRow
func (base *noDataBase) Stmt(stmt *sql.Stmt) (*sql.Stmt) {
	return nil
}






//------------------------------------




func (model *noDataModel) Create(data Map) (Map,error) {
	return nil,errors.New("无数据")
}
func (model *noDataModel) Change(item Map, data Map) (Map,error) {
	return nil,errors.New("无数据")
}
func (model *noDataModel) Remove(item Map) (error) {
	return nil
}
func (model *noDataModel) Entity(id Any) (Map,error) {
	return nil,errors.New("无数据")
}
func (model *noDataModel) Delete(args ...Map) (int64,error) {
	return int64(0),errors.New("无数据")
}
func (model *noDataModel) Update(args ...Map) (int64,error) {
	return int64(0),errors.New("无数据")
}
func (model *noDataModel) Count(args ...Map) (int64,error) {
	return int64(0),errors.New("无数据")
}
func (model *noDataModel) Single(args ...Map) (Map,error) {
	return nil,errors.New("无数据")
}
func (model *noDataModel) Query(args ...Map) ([]Map,error) {
	return nil,errors.New("无数据")
}
func (model *noDataModel) Limit(offset,limit Any, args ...Map) ([]Map,error) {
	return nil,errors.New("无数据")
}
func (model *noDataModel) Group(field string, args ...Map) ([]Map,error) {
	return nil,errors.New("无数据")
}