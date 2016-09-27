package base

import "time"


type (
	//触发器响应完成
	BodyTriggerFinish struct {}
	//触发器响应重来
	BodyTriggerRetrigger struct {
		Delay	time.Duration
	}
)



type (
	//计划响应完成
	BodyPlanFinish struct {}
	//计划响应重来
	BodyPlanReplan struct {
		Delay	time.Duration
	}
)

















type (

	//响应完成
	BodyFinish struct {

	}
	//响应重来
	BodyRetry struct {
		Delay	time.Duration
	}

	//响应文件下载
	BodyFile struct {
		//要下载的文件路径
		File string
		//自定义下载文件名
		Name string
	}
	//响应二进制下载
	BodyDown struct {
		//要下载的文件内容
		Body []byte
		//自定义下载文件名
		Name string
	}
	BodyJson struct {
		//响应的Json对象
		Json Any
	}
	BodyXml struct {
		//响应的Xml对象
		Xml Any
	}
	//响应跳转
	BodyGoto struct {
		Url string
	}
	//响应视图
	BodyView struct {
		View string
		Model Map
	}
)