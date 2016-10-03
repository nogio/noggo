package noggo



import (
	. "github.com/nogio/noggo/base"
)


type (
	//日志配置
	loggerConfig struct {
		//驱动
		Driver	string	`json:"driver"`
		//自定义配置
		Config	Map		`json:"config"`
	}
	//路由器配置
	routerConfig struct {
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


	//基础配置
	baseConfig struct {
		//计划驱动
		Driver	string	`json:"driver"`
		//计划驱动配置
		Config	Map		`json:"config"`

		//计划会话配置
		Session	*sessionConfig	`json:"session"`
		//计划路由器配置
		Router	*routerConfig	`json:"router"`
	}









	//触发器配置
	triggerConfig struct {
		//会话配置
		Session	*sessionConfig	`json:"session"`
		//路由器配置
		Router	*routerConfig	`json:"router"`
	}


	//计划配置
	planConfig struct {
		//计划驱动
		Driver	string	`json:"driver"`
		//计划驱动配置
		Config	Map		`json:"config"`

		//计划会话配置
		Session	*sessionConfig	`json:"session"`
		//计划路由器配置
		Router	*routerConfig	`json:"router"`
	}


	//任务配置
	taskConfig struct {
		//任务驱动
		Driver	string	`json:"driver"`
		//任务驱动配置
		Config	Map		`json:"config"`

		//任务会话配置
		Session	*sessionConfig	`json:"session"`
		//任务路由器配置
		Router	*routerConfig	`json:"router"`
	}


	//http配置
	httpConfig struct {
		//http驱动
		Driver	string	`json:"driver"`
		//http驱动配置
		Config	Map		`json:"config"`

		//http会话配置
		Session	*sessionConfig	`json:"session"`
		//http路由器配置
		Router	*routerConfig	`json:"router"`

		Port	string	`json:"port"`

		//字符集
		Charset	string	`json:"charset"`
	}





	//配置
	configConfig struct {
		Debug	bool			`json:"debug"`
		//默认路由配置
		Router	*routerConfig	`json:"router"`
		//默认会话配置
		Session	*sessionConfig	`json:"session"`

		//日志配置
		Logger	*loggerConfig	`json:"logger"`
		//触发器配置
		Trigger *triggerConfig	`json:"trigger"`

		//计划配置
		Plan *planConfig		`json:"plan"`
		//任务配置
		Task *taskConfig		`json:"task"`
		//http配置
		Http *httpConfig		`json:"http"`
	}

)
