package noggo

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	. "github.com/nogio/noggo/base"
	"errors"
	"os"
	"strings"
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
	//url
	Url *httpUrl
	//view
	View *viewGlobal

	//缓存
	Cache *cacheGlobal
	//数据
	Data *dataGlobal
	//存储
	Storage *storageGlobal

	//语法糖
	Sugar *sugarGlobal
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
		Http: &httpConfig{ Driver: "default", Config: Map{} },
		View: &viewConfig{ Driver: "default", Config: Map{} },
	}

	//常量模块
	Const = &constGlobal{}
	//Mapping模块
	Mapping = &mappingGlobal{}



	//日志模块
	Logger = &loggerGlobal{
		drivers: map[string]LoggerDriver{},
	}
	//会话模块
	Session = &sessionGlobal{
		drivers: map[string]SessionDriver{},
	}
	//触发器模块
	Trigger = &triggerGlobal{
		middlers:map[string]TriggerFunc{},
	}
	//任务模块
	Task = &taskGlobal{
		drivers: map[string]TaskDriver{},  middlers:map[string]TaskFunc{},
	}



	//计划模块
	Plan = &planGlobal{
		drivers: map[string]PlanDriver{}, middlers: map[string]PlanFunc{},
	}

	//事件模块
	Event = &eventGlobal{
		drivers: map[string]EventDriver{}, middlers: map[string]EventFunc{},
	}

	//队列模块
	Queue = &queueGlobal{
		drivers: map[string]QueueDriver{}, middlers: map[string]QueueFunc{},
		queueConnects: map[string]QueueConnect{},
	}
	//HTTP模块
	Http = &httpGlobal{
		drivers: map[string]HttpDriver{}, middlers:map[string]HttpFunc{},
	}

	Url = &httpUrl{}

	//View
	View = &viewGlobal{
		drivers: map[string]ViewDriver{}, helpers: map[string]Any{},
	}


	//缓存
	Cache = &cacheGlobal{
		drivers: map[string]CacheDriver{},
		connects: map[string]CacheConnect{},
	}
	//数据
	Data = &dataGlobal{
		drivers: map[string]DataDriver{},
		connects: map[string]DataConnect{},
	}
	//存储
	Storage = &storageGlobal{
		drivers: map[string]StorageDriver{},
		connects: map[string]StorageConnect{},
	}



	Sugar = &sugarGlobal{
		https: map[string]Map{},
	}

	Setting = Map{}


	loadConfig()
	loadLang()
}

//读取配置文件
func loadConfig() {
	err,cfg := readJsonConfig()
	if err != nil {
		panic("加载配置文件出错：" + err.Error())
	}
	Config = cfg

	//处理setting
	for k,v := range Config.Custom {
		Setting[k] = v
	}
	for k,v := range Config.Setting {
		Setting[k] = v
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

	configFile := "config.json"
	if len(os.Args) >= 2 {
		if strings.HasSuffix(os.Args[1],".json") {
			configFile = os.Args[1]	//第1个参数，0为程序本身
		}
	}


	bytes, err := ioutil.ReadFile(configFile)
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





