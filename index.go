package noggo

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"errors"
	"github.com/nogio/noggo/driver/http-default"
	"github.com/nogio/noggo/driver/plan-default"
	"github.com/nogio/noggo/driver/task-default"
	"github.com/nogio/noggo/driver/session-default"
	"github.com/nogio/noggo/driver/logger-default"
	"github.com/nogio/noggo/driver/view-default"
)

var (
	nodes map[string]*Noggo

	//当前位置
	//此变量比较重要，是用来注册对象时，表明当前的位置，如当前节点，当前数据库，等等
	Current string

	//全局配置
	Config *configConfig

	//常量模块
	Const *constGlobal
	//Map模块
	Mapping	*mappingGlobal




	//日志模块
	Logger *loggerGlobal

	//会话模块
	Session *sessionGlobal



	//触发器模块
	Trigger *triggerGlobal

	//任务模块
	Task *taskGlobal




	//计划模块
	Plan *planGlobal

	//http模块
	Http *httpGlobal


	//view
	View *viewGlobal



	//数据
	Data *dataGlobal

)



//框架初驱化
func init() {
	nodes = map[string]*Noggo{}
	//当前位置为空
	Current = ""

	//设置一个默认的配置
	Config = &configConfig{
		Debug: true,
		Lang: map[string]langConfig{},
		Node: map[string]nodeConfig{},

		Logger: &loggerConfig{ Driver: "default", Config: Map{} },
		Session: &sessionConfig{ Driver: "default", Config: Map{} },

		Trigger: &triggerConfig{ },
		Task: &taskConfig{ Driver: "default", Config: Map{} },

		Plan: &planConfig{ Driver: "default", Config: Map{} },
		Http: &httpConfig{ Driver: "default", Config: Map{}, Charset:"", Cookie:"noggo", Domain:"" },
		View: &viewConfig{ Driver: "default", Config: Map{} },
	}

	//常量模块
	Const = &constGlobal{}
	//Mapping模块
	Mapping = &mappingGlobal{}



	//日志模块
	Logger = &loggerGlobal{
		drivers: map[string]driver.LoggerDriver{},
	}
	//会话模块
	Session = &sessionGlobal{
		drivers: map[string]driver.SessionDriver{},
	}
	//触发器模块
	Trigger = &triggerGlobal{
		middlers:map[string]TriggerFunc{},
	}
	//任务模块
	Task = &taskGlobal{
		drivers: map[string]driver.TaskDriver{},  middlers:map[string]TaskFunc{},
	}



	//计划模块
	Plan = &planGlobal{
		drivers: map[string]driver.PlanDriver{}, middlers: map[string]PlanFunc{},
	}

	//HTTP模块
	Http = &httpGlobal{
		drivers: map[string]driver.HttpDriver{}, middlers:map[string]HttpFunc{},
	}

	//View
	View = &viewGlobal{
		drivers: map[string]driver.ViewDriver{}, helpers: map[string]Any{},
	}


	//数据
	Data = &dataGlobal{}

	loadConfig()
	loadLang()
	loadDriver()
}

//读取配置文件
func loadConfig() {
	err,cfg := readJsonConfig()
	if err == nil {
		Config = cfg
	}
}
func loadLang() {
	//加载多语言
	err, cfg := readJsonFile(fmt.Sprintf("langs/%v.json", ConstLangDefault))
	if err == nil {
		Const.Lang(ConstLangDefault, cfg)
	}
	for k,_ := range Config.Lang {
		err, cfg := readJsonFile(fmt.Sprintf("langs/%v.json", k))
		if err == nil {
			Const.Lang(k, cfg)
		}
	}
}
//加载默认驱动
func loadDriver() {
	Logger.Driver(ConstDriverDefault, logger_default.Driver())
	Session.Driver(ConstDriverDefault, session_default.Driver())

	Task.Driver(ConstDriverDefault, task_default.Driver())

	Plan.Driver(ConstDriverDefault, plan_default.Driver())
	Http.Driver(ConstDriverDefault, http_default.Driver())
	View.Driver(ConstDriverDefault, view_default.Driver())
}








func readJsonConfig() (error,*configConfig) {

	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		return err, nil
	} else {
		m := &configConfig{}
		err := json.Unmarshal(bytes, m)
		if err != nil {
			return err, nil
		} else {
			return nil, m
		}
	}
}

func readJsonFile(filename string) (error,Map) {
	bytes, err := ioutil.ReadFile(filename)
	if err == nil {
		m := make(Map)
		err := json.Unmarshal(bytes, &m)
		if err == nil {
			return nil, m
		}
	}
	return errors.New("读取失败"), nil
}







//-------------------------------- 语法糖 ----------------------------------------


//---------------------------------------------------------- 语法糖 begin ----------------------------------------------------------


//注册驱动
//方便注册各种驱动
func Driver(name string, drv Any) {
	switch v := drv.(type) {
	case driver.LoggerDriver:
		Logger.Driver(name, v)
	case driver.SessionDriver:
		Session.Driver(name, v)
	case driver.TaskDriver:
		Task.Driver(name, v)
	case driver.PlanDriver:
		Plan.Driver(name, v)
	case driver.HttpDriver:
		Http.Driver(name, v)
	case driver.ViewDriver:
		View.Driver(name, v)
	default:
		panic("未支持的驱动")
	}
}

//注册中间件
//方便注册各种中间件
func Middler(call Any) {
	switch v := call.(type) {

	case TriggerFunc:
		Trigger.Middler(NewMd5Id(), v)
	case func(*TriggerContext):
		Trigger.Middler(NewMd5Id(), v)

	case TaskFunc:
		Task.Middler(NewMd5Id(), v)
	case func(*TaskContext):
		Task.Middler(NewMd5Id(), v)

	case HttpFunc:
		Http.Middler(NewMd5Id(), v)
	case func(*HttpContext):
		Http.Middler(NewMd5Id(), v)

	case PlanFunc:
		Plan.Middler(NewMd5Id(), v)
	case func(*PlanContext):
		Plan.Middler(NewMd5Id(), v)

	default:
		panic("未支持的中间件")
	}
}
