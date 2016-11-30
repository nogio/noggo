package noggo

import (
	. "github.com/nogio/noggo/base"
	"sync"
	"errors"
	"database/sql"
	"fmt"
	"strings"
)

//data driver begin

const (
	ASC     = "$$$ASC$$$"
	DESC    = "$$$DESC$$$"
	COUNT   = "$$$COUNT$$$"
	AVG     = "$$$AVG$$$"
	SUM     = "$$$SUM$$$"
	MAX     = "$$$MAX$$$"
	MIN     = "$$$MIN$$$"
	DataFieldDelims  = "$$$"
)
type (
	FTS string  //全文搜索类型，解析sql的时候处理
	Nil struct {}
	NotNil struct {}
	DataDriver interface {
		//连接驱动的时候
		//应该做如下工作：
		//1. 检查config合法性
		//2. 初始化连接相关对象
		Connect(config Map) (DataConnect,error)
	}
	DataConnect interface {
		//打开数据库连接
		Open() (error)
		//关闭数据库连接
		Close() (error)

		//构建模型和索引
		Build() error

		//注册模型
		Model(string,Map)
		//注册视图
		View(string,Map)

		//获取数据库对象
		Base(string,CacheBase) (DataBase,error)
	}

	//数据库接口
	DataBase interface {
		Close()
		Model(name string) (DataModel)
		View(name string) (DataView)

		//是否使用缓存，默认使用
		Cache(bool) (DataBase)
		//开启手动提交事务模式
		Begin() (DataBase)
		Submit() (error)
		Cancel() (error)

		//原生SQL的方法，接口可以执行原生查询，以支持Model不能完成的工作
		//但是，必须 调用Begin之后，才能使用下例方法，然后 Submit 或 Cancel
		//因为全部使用事务。
		Exec(query string, args ...interface{}) (sql.Result, error)
		Prepare(query string) (*sql.Stmt, error)
		Query(query string, args ...interface{}) (*sql.Rows, error)
		QueryRow(query string, args ...interface{}) (*sql.Row)
		Stmt(stmt *sql.Stmt) (*sql.Stmt)
	}


	//数据视图接口
	DataView interface {
		Count(args ...Any) (int64,error)
		Single(args ...Any) (Map,error)
		Query(args ...Any) ([]Map,error)
		Limit(offset, limit Any, args ...Any) (int64,[]Map,error)
		Group(field string, args ...Any) ([]Map,error)
	}

	//数据模型接口
	DataModel interface {
		DataView

		Create(Map) (Map,error)
		Change(Map,Map) (Map,error)
		Remove(Map) (error)
		Entity(Any) (Map,error)

		Update(sets Map,args ...Any) (int64,error)
		Delete(args ...Any) (int64,error)
	}

)
//data driver end




type (
	dataGlobal    struct {
		mutex       sync.Mutex
		drivers     map[string]DataDriver

		//数据连接
		connects    map[string]DataConnect

		models      map[string]map[string]Map
		views      map[string]map[string]Map
	}
)


//连接驱动
func (global *dataGlobal) connect(config *dataConfig) (DataConnect,error) {
	if dataDriver,ok := global.drivers[config.Driver]; ok {
		return dataDriver.Connect(config.Config)
	} else {
		panic("数据：不支持的驱动 " + config.Driver)
	}
}


//注册数据驱动
func (global *dataGlobal) Driver(name string, config DataDriver) {
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


				//注册视图
				//先全局
				for k,v := range global.views[ConstNodeGlobal] {
					conn.View(k, v)
				}
				//再库
				for k,v := range global.views[name] {
					conn.View(k, v)
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


//注册视图
func (global *dataGlobal) View(name string, config Map) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.views == nil {
		global.views = map[string]map[string]Map{}
	}

	//节点
	nodeName := ConstNodeGlobal
	if Current != "" {
		nodeName = Current
	}

	//如果节点配置不存在，创建
	if global.views[nodeName] == nil {
		global.views[nodeName] = map[string]Map{}
	}

	//可以后注册重写原有路由配置，所以直接保存
	global.views[nodeName][name] = config
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


//取模型的枚举定义
func (global *dataGlobal) Enums(data, model, field string) (Map) {

	if Data.models[data] != nil {
		dataConfig := Data.models[data]
		if dataConfig[model] != nil {
			modelConfig := dataConfig[model]
			if modelConfig["fields"] != nil {
				fields := modelConfig["fields"].(Map)
				if fields[field] != nil {
					fieldConfig := fields[field].(Map)
					if fieldConfig["enum"] != nil {
						return fieldConfig["enum"].(Map)
					}
				}
			}
		}
	}

	return Map{}
}





































//返回DB对象
func (global *dataGlobal) Base(name string) (DataBase) {
	if conn,ok := global.connects[name]; ok {

		//缓存相关
		var cb CacheBase = nil
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



//查询语法解析器
// 字段包裹成  $$$field$$$ 请自行处理
// 如mysql为反引号`field`，postgres为引号"field"，
// 所有参数使问号(?)
// postgres驱动需要自行处理转成 $1,$2这样的
func (global *dataGlobal) Parsing(args ...Any) (string,[]interface{},string,error) {

	if len(args) > 0 {

		//如果直接写sql
		if v,ok := args[0].(string); ok {
			sql := v
			params := []interface{}{}

			for i,arg := range args {
				if i > 0 {
					params = append(params, arg)
				}
			}

			return sql,params,"",nil

		} else {


			maps := []Map{}
			for _,v := range args {
				if m,ok := v.(Map); ok {
					maps = append(maps, m)
				}
			}

			querys,values,orders := global.parsing(maps...)

			orderStr := ""
			if len(orders) > 0 {
				orderStr = fmt.Sprintf("ORDER BY %s", strings.Join(orders, ","))
			}

			//sql := fmt.Sprintf("%s %s", strings.Join(querys, " OR "), orderStr)

			return strings.Join(querys, " OR "), values, orderStr, nil
		}
	} else {
		return "1=1",[]interface{}{},"",nil
	}
}
//注意，这个是实际的解析，支持递归
func (global *dataGlobal) parsing(args ...Map) ([]string,[]interface{},[]string) {

	fp := DataFieldDelims

	querys := []string{}
	values := make([]interface{}, 0)
	orders := []string{}

	//否则是多个map,单个为 与, 多个为 或
	for _,m := range args {
		ands := []string{}

		for k,v := range m {

			//如果值是ASC,DESC，表示是排序
			if ov,ok := v.(string); ok && (ov==ASC || ov==DESC) {

				if ov == ASC {
					orders = append(orders, fmt.Sprintf(`%s%s%s ASC`, fp, k, fp))
				} else {
					orders = append(orders, fmt.Sprintf(`%s%s%s DESC`, fp, k, fp))
				}

			} else if ms,ok := v.([]Map); ok {
				//是[]Map，相当于or

				qs,vs,os := global.parsing(ms...)
				if len(qs) > 0 {
					ands = append(ands, fmt.Sprintf("(%s)", strings.Join(qs, " OR ")))
					for _,vsVal := range vs {
						values = append(values, vsVal)
					}
				}
				for _,osVal := range os {
					orders = append(orders, osVal)
				}

			} else {

				//v要处理一下如果是map要特别处理
				//key做为操作符，比如 > < >= 等
				//而且多个条件是and，比如 views > 1 AND views < 100
				//自定义操作符的时候，可以用  is not null 吗？
				if opMap, opOK := v.(Map); opOK {

					opAnds := []string{}
					for opKey,opVal := range opMap {
						opAnds = append(opAnds, fmt.Sprintf(`%s%s%s %s ?`, fp, k, fp, opKey))
						values = append(values, opVal)
					}
					ands = append(ands, fmt.Sprintf("(%s)", strings.Join(opAnds, " AND ")))

				} else {

					if v == nil {
						ands = append(ands, fmt.Sprintf(`%s%s%s IS NULL`, fp, k, fp))
					} else if _,ok := v.(Nil); ok {
						//为空值
						ands = append(ands, fmt.Sprintf(`%s%s%s IS NULL`, fp, k, fp))
					} else if _,ok := v.(NotNil); ok {
						//不为空值
						ands = append(ands, fmt.Sprintf(`%s%s%s IS NOT NULL`, fp, k, fp))
					} else if fts,ok := v.(FTS); ok {
						//处理模糊搜索
						safeFts := strings.Replace(string(fts), "'", "''", -1)
						ands = append(ands, fmt.Sprintf(`%s%s%s LIKE '%%%s%%'`, fp, k, fp, safeFts))
					} else {
						ands = append(ands, fmt.Sprintf(`%s%s%s = ?`, fp, k, fp))
						values = append(values, v)
					}
				}

			}

		}

		if len(ands) > 0 {
			querys = append(querys, fmt.Sprintf("(%s)", strings.Join(ands, " AND ")))
		}
	}

	return querys,values,orders
}






type (
	noDataBase struct {}
	noDataModel struct {
		noDataView
	}
	noDataView struct {}
)
func (base *noDataBase) Close() {
}
func (base *noDataBase) Model(name string) (DataModel) {
	return &noDataModel{}
}
func (base *noDataBase) View(name string) (DataView) {
	return &noDataView{}
}
func (base *noDataBase) Cache(use bool)(DataBase) {
	return base
}
func (base *noDataBase) Begin() (DataBase) {
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
func (model *noDataModel) Delete(args ...Any) (int64,error) {
	return int64(0),errors.New("无数据")
}
func (model *noDataModel) Update(sets Map,args ...Any) (int64,error) {
	return int64(0),errors.New("无数据")
}
/*
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
*/




func (model *noDataView) Count(args ...Any) (int64,error) {
	return int64(0),errors.New("无数据")
}
func (model *noDataView) Single(args ...Any) (Map,error) {
	return Map{},errors.New("无数据")
}
func (model *noDataView) Query(args ...Any) ([]Map,error) {
	return []Map{},errors.New("无数据")
}
func (model *noDataView) Limit(offset,limit Any, args ...Any) (int64,[]Map,error) {
	return int64(0),[]Map{},errors.New("无数据")
}
func (model *noDataView) Group(field string, args ...Any) ([]Map,error) {
	return []Map{},errors.New("无数据")
}