package noggo



import (
	. "github.com/nogio/noggo/base"
)


type (
	//语言配置
	langConfig struct {
		Name		string		`json:"name"`
		Text		string		`json:"Text"`
		Accepts		[]string	`json:"accepts"`
	}

	//节点配置
	nodeConfig struct {
		Id	string	`json:"id"`
		Port string	`json:"port"`
		Url string	`json:"url"`
		Name string	`json:"name"`
		Text string	`json:"text"`

		Crypto string `json:"crypto"`
		//字符集
		Charset	string	`json:"charset"`
		//Session使用的cookie name
		Cookie	string	`json:"cookie"`
		//Session使用的domain
		Domain	string	`json:"domain"`

		//节点指定SESSION
		Session *sessionConfig `json:"session"`
	}

	//data配置
	dataConfig struct {
		//data驱动
		Driver	string	`json:"driver"`
		//data驱动配置
		Config	Map		`json:"config"`
		//缓存key
		Cache	string	`json:"cache"`
	}

	//cache配置
	cacheConfig struct {
		//驱动
		Driver	string	`json:"driver"`
		//驱动配置
		Config	Map		`json:"config"`
	}



	//日志配置
	loggerConfig struct {
		//驱动
		Driver	string	`json:"driver"`
		//自定义配置
		Config	Map		`json:"config"`
	}
	//会话配置
	sessionConfig struct {
		//驱动
		Driver	string	`json:"driver"`
		//超时时间
		Expiry	int64	`json:"expiry"`
		//自定义配置
		Config	Map		`json:"config"`
	}
	//触发器配置
	triggerConfig struct {

	}

	//任务配置
	taskConfig struct {
		//任务驱动
		Driver	string	`json:"driver"`
		//任务驱动配置
		Config	Map		`json:"config"`
	}







	//计划配置
	planConfig struct {
		//计划驱动
		Driver	string	`json:"driver"`
		//计划驱动配置
		Config	Map		`json:"config"`
	}



	//事件配置
	eventConfig struct {
		//驱动
		Driver	string	`json:"driver"`
		//前缀
		Prefix	string	`json:"prefix"`
		//驱动配置
		Config	Map		`json:"config"`
	}


	//queue配置
	queueConfig struct {
		//驱动
		Driver	string	`json:"driver"`
		//前缀
		Prefix	string	`json:"prefix"`
		//驱动配置
		Config	Map		`json:"config"`
	}




	//http配置
	httpConfig struct {
		//http驱动
		Driver	string	`json:"driver"`
		//http驱动配置
		Config	Map		`json:"config"`
	}





	//view配置
	viewConfig struct {
		//view驱动
		Driver	string	`json:"driver"`
		//view驱动配置
		Config	Map		`json:"config"`
	}


	//配置
	configConfig struct {
		Debug	bool			`json:"debug"`

		Lang	map[string]*langConfig	`json:"lang"`
		Node	map[string]*nodeConfig	`json:"node"`
		Data	map[string]*dataConfig	`json:"data"`
		Cache	map[string]*cacheConfig	`json:"cache"`
		Queue	map[string]*queueConfig	`json:"queue"`


		//日志配置
		Logger	*loggerConfig	`json:"logger"`
		//默认会话配置
		Session	*sessionConfig	`json:"session"`
		//触发器配置
		Trigger *triggerConfig	`json:"trigger"`
		//任务配置
		Task *taskConfig		`json:"task"`



		//计划配置
		Plan *planConfig		`json:"plan"`
		//事件配置
		Event *eventConfig		`json:"event"`



		//http配置
		Http *httpConfig		`json:"http"`

		//view
		View *viewConfig		`json:"view"`


		//自定义配置
		Custom Map		`json:"custom"`
	}

)
