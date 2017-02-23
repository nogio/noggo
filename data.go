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
	/*
	//已经废弃，改到base包中了，暂时存留for兼容
	ASC     = "$$$ASC$$$"
	DESC    = "$$$DESC$$$"
	COUNT   = "$$$COUNT$$$"
	AVG     = "$$$AVG$$$"
	SUM     = "$$$SUM$$$"
	MAX     = "$$$MAX$$$"
	MIN     = "$$$MIN$$$"
	DataFieldDelims  = "$"
	*/
)
type (
	/*
	//已经废弃，改到base包中，暂时存留for兼容
	FTS string  //全文搜索类型，解析sql的时候处理
	Nil struct {}
	NotNil struct {}
	*/







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

		//注册表
		Table(string,Map)
		//注册视图
		View(string,Map)
		//注册模型
		Model(string,Map)

		//获取数据库对象
		Base(string,CacheBase) (DataBase,error)
	}
	DataTrigger struct {
		Name    string
		Value   Map
	}
	//数据库接口
	DataBase interface {
		Close()
		Table(name string) (DataTable)
		View(name string) (DataView)
		Model(name string) (DataModel)

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

	//数据表接口
	DataTable interface {
		Create(Map) (Map,error)
		Change(Map,Map) (Map,error)
		Remove(Map) (error)
		Recover(Map) (error)

		Update(sets Map,args ...Any) (int64,error)
		Delete(args ...Any) (int64,error)

		Entity(Any) (Map,error)
		Count(args ...Any) (int64,error)
		Single(args ...Any) (Map,error)
		Query(args ...Any) ([]Map,error)
		//Querys(keyword string, args ...Any) ([]Map,error)
		Limit(offset, limit Any, args ...Any) (int64,[]Map,error)
		//Limits(offset, limit Any, keyword string, args ...Any) (int64,[]Map,error)
		Group(field string, args ...Any) ([]Map,error)

	}

	//数据视图接口
	DataView interface {
		Count(args ...Any) (int64,error)
		Single(args ...Any) (Map,error)
		Query(args ...Any) ([]Map,error)
		Limit(offset, limit Any, args ...Any) (int64,[]Map,error)
		Group(field string, args ...Any) ([]Map,error)
		Entity(Any) (Map,error)
	}

	//数据模型接口
	DataModel interface {
		Single(args ...Any) (Map,error)
		Query(args ...Any) ([]Map,error)
	}

)
//data driver end




type (
	dataGlobal    struct {
		mutex       sync.Mutex
		drivers     map[string]DataDriver

		//数据连接
		connects    map[string]DataConnect

		//三大对象：表，视图，模型
		tables      map[string]map[string]Map
		views      map[string]map[string]Map
		models      map[string]map[string]Map
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

				//注册表
				//先全局
				for k,v := range global.tables[ConstNodeGlobal] {
					conn.Table(k, v)
				}
				//再库
				for k,v := range global.tables[name] {
					conn.Table(k, v)
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


				//注册模型
				//先全局
				for k,v := range global.models[ConstNodeGlobal] {
					conn.Model(k, v)
				}
				//再库
				for k,v := range global.models[name] {
					conn.Model(k, v)
				}


				//如果需要构建表和索引
				if config.Build {
					conn.Build()
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






















//注册表
func (global *dataGlobal) Table(name string, config Map) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.tables == nil {
		global.tables = map[string]map[string]Map{}
	}

	//节点
	nodeName := ConstNodeGlobal
	if Current != "" {
		nodeName = Current
	}

	//如果节点配置不存在，创建
	if global.tables[nodeName] == nil {
		global.tables[nodeName] = map[string]Map{}
	}

	//可以后注册重写原有配置，所以直接保存
	global.tables[nodeName][name] = config
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


	//可以后注册重写原有配置，所以直接保存
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

	//可以后注册重写原有配置，所以直接保存
	global.views[nodeName][name] = config
}













//查询某库所有表
func (global *dataGlobal) Tables(bases ...string) (Map) {

	if len(bases) > 0 {
		base := bases[0]
		if global.tables[base] != nil {

			m := Map{}
			for k,v := range global.tables[base] {
				m[k] = v
			}
			return m
		}
	} else {
		m := Map{}
		for k,v := range global.tables {
			m[k] = v
		}
		return m
	}

	return Map{}
}

//查询某库所有视图
func (global *dataGlobal) Views(bases ...string) (Map) {

	if len(bases) > 0 {
		base := bases[0]
		if global.views[base] != nil {

			m := Map{}
			for k,v := range global.views[base] {
				m[k] = v
			}
			return m
		}
	} else {
		m := Map{}
		for k,v := range global.views {
			m[k] = v
		}
		return m
	}

	return Map{}
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














//查询某表部分字段
func (global *dataGlobal) TableField(base, name string, fields []string, maps ...Map) Map {

	m := Map{}

	if dc,ok := global.tables[base]; ok {
		if mc,ok := dc[name]; ok {
			if fc,ok := mc["field"]; ok {

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
//查询某表所有字段
func (global *dataGlobal) TableFields(base, name string, maps ...Map) (Map) {

	m := Map{}

	if dc,ok := global.tables[base]; ok {
		if mc,ok := dc[name]; ok {
			if fc,ok := mc["field"]; ok {
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
//取模型的枚举定义
func (global *dataGlobal) TableEnums(data, name, field string) (Map) {

	if dataConfig,ok := global.tables[data]; ok {
		if tableConfig,ok := dataConfig[name]; ok {
			if fields,ok := tableConfig["field"].(Map); ok {
				if fieldConfig,ok := fields[field].(Map); ok {
					if vv,ok := fieldConfig["enum"].(Map); ok {
						return vv
					}
				}
			}
		}
	}
	return Map{}
}































//查询某视图部分字段
func (global *dataGlobal) ViewField(base, name string, fields []string, maps ...Map) Map {

	m := Map{}

	if dc,ok := global.views[base]; ok {
		if mc,ok := dc[name]; ok {
			if fc,ok := mc["field"]; ok {

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
//查询某视图部分字段
func (global *dataGlobal) ViewFields(base, name string, maps ...Map) (Map) {

	m := Map{}

	if dc,ok := global.views[base]; ok {
		if mc,ok := dc[name]; ok {
			if fc,ok := mc["field"]; ok {
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
//查询某视图部分字段
func (global *dataGlobal) ViewEnums(data, name, field string) (Map) {

	if dataConfig,ok := global.views[data]; ok {
		if tableConfig,ok := dataConfig[name]; ok {
			if fields,ok := tableConfig["field"].(Map); ok {
				if fieldConfig,ok := fields[field].(Map); ok {
					if vv,ok := fieldConfig["enum"].(Map); ok {
						return vv
					}
				}
			}
		}
	}
	return Map{}
}









//查询某视图部分字段
func (global *dataGlobal) ModelField(base, name string, fields []string, maps ...Map) Map {

	m := Map{}

	if dc,ok := global.models[base]; ok {
		if mc,ok := dc[name]; ok {
			if fc,ok := mc["field"]; ok {

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
//查询某视图部分字段
func (global *dataGlobal) ModelFields(base, name string, maps ...Map) (Map) {

	m := Map{}

	if dc,ok := global.models[base]; ok {
		if mc,ok := dc[name]; ok {
			if fc,ok := mc["field"]; ok {
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
//查询某视图部分字段
func (global *dataGlobal) ModelEnums(data, name, field string) (Map) {

	if dataConfig,ok := global.models[data]; ok {
		if tableConfig,ok := dataConfig[name]; ok {
			if fields,ok := tableConfig["field"].(Map); ok {
				if fieldConfig,ok := fields[field].(Map); ok {
					if vv,ok := fieldConfig["enum"].(Map); ok {
						return vv
					}
				}
			}
		}
	}
	return Map{}
}


















//查询部分字段，跨三个对象
func (global *dataGlobal) Field(base, name string, fields []string, maps ...Map) Map {

	m := Map{}

	if dc,ok := global.models[base]; ok {
		if mc,ok := dc[name]; ok {
			if fc,ok := mc["field"]; ok {

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
	if dc,ok := global.views[base]; ok {
		if mc,ok := dc[name]; ok {
			if fc,ok := mc["field"]; ok {

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
	if dc,ok := global.tables[base]; ok {
		if mc,ok := dc[name]; ok {
			if fc,ok := mc["field"]; ok {

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
//查询某字段，跨三个对象
func (global *dataGlobal) Fields(base, name string, maps ...Map) (Map) {

	m := Map{}

	if dc,ok := global.models[base]; ok {
		if mc,ok := dc[name]; ok {
			if fc,ok := mc["field"]; ok {
				//不可以直接给,要给一个新的,要不么返回了引用,改了后, 原定义也改了
				for k,v := range fc.(Map) {
					m[k] = v
				}
			}
		}
	}
	if dc,ok := global.views[base]; ok {
		if mc,ok := dc[name]; ok {
			if fc,ok := mc["field"]; ok {
				//不可以直接给,要给一个新的,要不么返回了引用,改了后, 原定义也改了
				for k,v := range fc.(Map) {
					m[k] = v
				}
			}
		}
	}
	if dc,ok := global.tables[base]; ok {
		if mc,ok := dc[name]; ok {
			if fc,ok := mc["field"]; ok {
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
//查询某枚举，跨三个对象
func (global *dataGlobal) Enums(data, name, field string) (Map) {

	m := Map{}

	if dataConfig,ok := global.models[data]; ok {
		if tableConfig,ok := dataConfig[name]; ok {
			if fields,ok := tableConfig["field"].(Map); ok {
				if fieldConfig,ok := fields[field].(Map); ok {
					if vv,ok := fieldConfig["enum"].(Map); ok {

						for k,v := range vv {
							m[k] = v
						}

					}
				}
			}
		}
	}

	if dataConfig,ok := global.views[data]; ok {
		if tableConfig,ok := dataConfig[name]; ok {
			if fields,ok := tableConfig["field"].(Map); ok {
				if fieldConfig,ok := fields[field].(Map); ok {
					if vv,ok := fieldConfig["enum"].(Map); ok {

						for k,v := range vv {
							m[k] = v
						}

					}
				}
			}
		}
	}

	if dataConfig,ok := global.tables[data]; ok {
		if tableConfig,ok := dataConfig[name]; ok {
			if fields,ok := tableConfig["field"].(Map); ok {
				if fieldConfig,ok := fields[field].(Map); ok {
					if vv,ok := fieldConfig["enum"].(Map); ok {

						for k,v := range vv {
							m[k] = v
						}

					}
				}
			}
		}
	}

	return m
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
// 字段包裹成  $field$ 请自行处理
// 如mysql为反引号`field`，postgres, oracle为引号"field"，
// 所有参数使问号(?)
// postgres驱动需要自行处理转成 $1,$2这样的
// oracle驱动需要自行处理转成 :1 :2这样的
func (global *dataGlobal) Parsing(args ...Any) (string,[]interface{},string,error) {

	if len(args) > 0 {

		//如果直接写sql
		if v,ok := args[0].(string); ok {
			sql := v
			params := []interface{}{}
			orderBy := ""

			for i,arg := range args {
				if i > 0 {
					params = append(params, arg)
				}
			}

			//这里要处理一下，把order提取出来
			//先拿到 order by 的位置
			i := strings.Index(strings.ToLower(sql), "order by")
			if i >= 0 {
				orderBy = sql[i:]
				sql = sql[:i]
			}

			return sql,params,orderBy,nil

		} else {


			maps := []Map{}
			for _,v := range args {
				if m,ok := v.(Map); ok {
					maps = append(maps, m)
				}
				//如果直接是[]Map，应该算OR处理啊，暂不处理这个
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

	fp := DELIMS

	querys := []string{}
	values := make([]interface{}, 0)
	orders := []string{}

	//否则是多个map,单个为 与, 多个为 或
	for _,m := range args {
		ands := []string{}

		for k,v := range m {

			//如果值是ASC,DESC，表示是排序
			//if ov,ok := v.(string); ok && (ov==ASC || ov==DESC) {
			if v==ASC {
				//正序
				orders = append(orders, fmt.Sprintf(`%s%s%s ASC`, fp, k, fp))
			} else if v==DESC {
				//倒序
				orders = append(orders, fmt.Sprintf(`%s%s%s DESC`, fp, k, fp))

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
						//这里要支持LIKE
						if opKey == LIKE {
							safeFts := strings.Replace(fmt.Sprintf("%v", opVal), "'", "''", -1)
							opAnds = append(opAnds, fmt.Sprintf(`%s%s%s LIKE '%%%s%%'`, fp, k, fp, safeFts))
						} else if opKey == FULL {
							safeFts := strings.Replace(fmt.Sprintf("%v", opVal), "'", "''", -1)
							opAnds = append(opAnds, fmt.Sprintf(`%s%s%s LIKE '%%%s%%'`, fp, k, fp, safeFts))
						} else if opKey == LEFT {
							safeFts := strings.Replace(fmt.Sprintf("%v", opVal), "'", "''", -1)
							opAnds = append(opAnds, fmt.Sprintf(`%s%s%s LIKE '%s%%'`, fp, k, fp, safeFts))
						} else if opKey == RIGHT {
							safeFts := strings.Replace(fmt.Sprintf("%v", opVal), "'", "''", -1)
							opAnds = append(opAnds, fmt.Sprintf(`%s%s%s LIKE '%%%s'`, fp, k, fp, safeFts))
						} else if opKey == IN {
							//IN (?,?,?)

							//safeFts := strings.Replace(fmt.Sprintf("%v", opVal), "'", "''", -1)
							//opAnds = append(opAnds, fmt.Sprintf(`%s%s%s LIKE '%%%s'`, fp, k, fp, safeFts))


						} else {
							opAnds = append(opAnds, fmt.Sprintf(`%s%s%s %s ?`, fp, k, fp, opKey))
							values = append(values, opVal)
						}
					}


					ands = append(ands, fmt.Sprintf("(%s)", strings.Join(opAnds, " AND ")))

				} else {

					if v == nil {
						ands = append(ands, fmt.Sprintf(`%s%s%s IS NULL`, fp, k, fp))
					} else if v == NIL {
						ands = append(ands, fmt.Sprintf(`%s%s%s IS NULL`, fp, k, fp))
					} else if v == NOL {
						//不为空值
						ands = append(ands, fmt.Sprintf(`%s%s%s IS NOT NULL`, fp, k, fp))
						/*
					}  else if _,ok := v.(Nil); ok {
						//为空值
						ands = append(ands, fmt.Sprintf(`%s%s%s IS NULL`, fp, k, fp))
					} else if _,ok := v.(NotNil); ok {
						//不为空值
						ands = append(ands, fmt.Sprintf(`%s%s%s IS NOT NULL`, fp, k, fp))
					} else if fts,ok := v.(FTS); ok {
						//处理模糊搜索，此条后续版本会移除
						safeFts := strings.Replace(string(fts), "'", "''", -1)
						ands = append(ands, fmt.Sprintf(`%s%s%s LIKE '%%%s%%'`, fp, k, fp, safeFts))
						*/
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
	noDataTable struct {}
	noDataView struct {}
	noDataModel struct {}
)
func (base *noDataBase) Close() {
}
func (base *noDataBase) Table(name string) (DataTable) {
	return &noDataTable{}
}
func (base *noDataBase) View(name string) (DataView) {
	return &noDataView{}
}
func (base *noDataBase) Model(name string) (DataModel) {
	return &noDataModel{}
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




func (table *noDataTable) Create(data Map) (Map,error) {
	return nil,errors.New("无数据")
}
func (table *noDataTable) Change(item Map, data Map) (Map,error) {
	return nil,errors.New("无数据")
}
func (table *noDataTable) Remove(item Map) (error) {
	return nil
}
func (table *noDataTable) Recover(item Map) (error) {
	return nil
}
func (table *noDataTable) Delete(args ...Any) (int64,error) {
	return int64(0),errors.New("无数据")
}
func (table *noDataTable) Update(sets Map,args ...Any) (int64,error) {
	return int64(0),errors.New("无数据")
}


func (view *noDataTable) Count(args ...Any) (int64,error) {
	return int64(0),errors.New("无数据")
}
func (view *noDataTable) Single(args ...Any) (Map,error) {
	return Map{},errors.New("无数据")
}
func (view *noDataTable) Query(args ...Any) ([]Map,error) {
	return []Map{},errors.New("无数据")
}
func (view *noDataTable) Querys(keyword string,args ...Any) ([]Map,error) {
	return []Map{},errors.New("无数据")
}
func (view *noDataTable) Limit(offset,limit Any, args ...Any) (int64,[]Map,error) {
	return int64(0),[]Map{},errors.New("无数据")
}
func (view *noDataTable) Limits(offset,limit Any, keyword string, args ...Any) (int64,[]Map,error) {
	return int64(0),[]Map{},errors.New("无数据")
}
func (view *noDataTable) Group(field string, args ...Any) ([]Map,error) {
	return []Map{},errors.New("无数据")
}
func (view *noDataTable) Entity(id Any) (Map,error) {
	return nil,errors.New("无数据")
}





func (view *noDataView) Count(args ...Any) (int64,error) {
	return int64(0),errors.New("无数据")
}
func (view *noDataView) Single(args ...Any) (Map,error) {
	return Map{},errors.New("无数据")
}
func (view *noDataView) Query(args ...Any) ([]Map,error) {
	return []Map{},errors.New("无数据")
}
func (view *noDataView) Limit(offset,limit Any, args ...Any) (int64,[]Map,error) {
	return int64(0),[]Map{},errors.New("无数据")
}
func (view *noDataView) Group(field string, args ...Any) ([]Map,error) {
	return []Map{},errors.New("无数据")
}
func (view *noDataView) Entity(id Any) (Map,error) {
	return nil,errors.New("无数据")
}








func (view *noDataModel) Single(args ...Any) (Map,error) {
	return Map{},errors.New("无数据")
}
func (view *noDataModel) Query(args ...Any) ([]Map,error) {
	return []Map{},errors.New("无数据")
}