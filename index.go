package noggo

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"errors"
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

	//事件模块
	Event *eventGlobal

	//队列模块
	Queue *queueGlobal

	//http模块
	Http *httpGlobal


	//view
	View *viewGlobal



	//缓存
	Cache *cacheGlobal

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
		Lang: map[string]*langConfig{},
		Node: map[string]*nodeConfig{},
		Data: map[string]*dataConfig{},

		Logger: &loggerConfig{ Driver: "default", Config: Map{} },
		Session: &sessionConfig{ Driver: "default", Config: Map{} },

		Trigger: &triggerConfig{ },
		Task: &taskConfig{ Driver: "default", Config: Map{} },

		Plan: &planConfig{ Driver: "default", Config: Map{} },
		Event: &eventConfig{ Driver: "default", Prefix: "", Config: Map{} },
		Queue: map[string]*queueConfig{ "default": &queueConfig{ Driver: "default", Prefix: "", Config: Map{} } },
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

	//事件模块
	Event = &eventGlobal{
		drivers: map[string]driver.EventDriver{}, middlers: map[string]EventFunc{},
	}

	//队列模块
	Queue = &queueGlobal{
		drivers: map[string]driver.QueueDriver{}, middlers: map[string]QueueFunc{},
		queueConnects: map[string]driver.QueueConnect{},
	}
	//HTTP模块
	Http = &httpGlobal{
		drivers: map[string]driver.HttpDriver{}, middlers:map[string]HttpFunc{},
	}

	//View
	View = &viewGlobal{
		drivers: map[string]driver.ViewDriver{}, helpers: map[string]Any{},
	}


	//缓存
	Cache = &cacheGlobal{
		drivers: map[string]driver.CacheDriver{},
		connects: map[string]driver.CacheConnect{},
	}
	//数据
	Data = &dataGlobal{
		drivers: map[string]driver.DataDriver{},
		connects: map[string]driver.DataConnect{},
	}

	loadConfig()
	loadLang()
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
		Const.Langs(ConstLangDefault, cfg)
	}
	for k,_ := range Config.Lang {
		err, cfg := readJsonFile(fmt.Sprintf("langs/%v.json", k))
		if err == nil {
			Const.Langs(k, cfg)
		}
	}
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

	case driver.EventDriver:
		Event.Driver(name, v)

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


	case EventFunc:
		Event.Middler(NewMd5Id(), v)
	case func(*EventContext):
		Event.Middler(NewMd5Id(), v)

	default:
		panic("未支持的中间件")
	}
}
