package noggo

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	. "github.com/nogio/noggo/base"
	"errors"
)

var (
	nodes map[string]*Noggo

	//当前位置
	//此变量比较重要，是用来注册对象时，表明当前的位置，如当前节点，当前数据库，等等
	Current string

	//全局配置
	Config configConfig

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

	/*


	//路由器模块
	Router *routerModule


	//任务模块
	Task *taskModule
	*/
)



//框架初驱化
func init() {
	nodes = map[string]*Noggo{}
	//当前位置为空
	Current = ""
	//读取配置文件
	Config = readJsonConfig()

	//常量模块
	Const = &constGlobal{}
	//Mapping模块
	Mapping = &mappingGlobal{}
	//日志模块
	Logger = &loggerGlobal{}
	//会话模块
	Session = &sessionGlobal{}
	//触发器模块
	Trigger = &triggerGlobal{}
	//任务模块
	Task = &taskGlobal{}



	//计划模块
	Plan = &planGlobal{}

	//HTTP模块
	Http = &httpGlobal{}



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


func readJsonConfig() (configConfig) {
	m := configConfig{}

	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(fmt.Sprintf("读取配置文件出错:%v", err))
	} else {
		err := json.Unmarshal(bytes, &m)
		if err != nil {
			panic(fmt.Sprintf("解析配置文件出错:%v", err))
		}
	}

	return m
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