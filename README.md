# 注意

此框架刚完成，尚未正式发布，目前我们正在一些正式项目中使用此框架，以便对框架进行细节的打磨和完善，请暂时不要在生产环境使用此框架，可能还会有很多调整和优化的地方。
20161104

# 简介

noggo 一个可以直接拿来干活的开发框架。

# 目标

noggo所追求的是开发效率，而非单节点的极致性能，若要提升性能可以有多种手段实现。
noggo整合后端开发常用的，日志、触发器、任务、计划、事件(pub/sub)、队列，数据库（SQL,NOSQL)、缓存、存储等模块，并提供一致性的开发体验，降低后端开发门槛，提升后端开发效率。
noggo可以让您使用一套框架，完成整个后端的开发工作，并提供自动生成文档等功能，为后端开发与前端开发的对接提供无限便利。


# 联络

官方QQ群：34613300

# 特色

*   模块驱动化
*   中间件
*   多语言支持
*   自动生成文档
*   一致性开发体验



# 示例

下载代码：

```
go get github.com/nogio/noggo
```

基础示例：

```golang
package main

import (
	"github.com/nogio/noggo"
	_ "github.com/nogio/noggo/core"
	"github.com/nogio/noggo/middler"
)


func main() {

	//请求日志、静态文件、表单解析、中间件
	noggo.Use(middler.HttpLogger())
	noggo.Use(middler.HttpStatic("statics"))
	noggo.Use(middler.HttpForm("uploads"))

	//get 首页
	noggo.Get("/", func(ctx *noggo.HttpContext) {
		ctx.Text("hello noggo.")
	})
	//post 首页
	noggo.Post("/", func(ctx *noggo.HttpContext) {
		ctx.Json(ctx.Form)
	})
	//添加一个触发器，noggo.Trigger.Touch("trigger.test") 触发
	noggo.Add("trigger.test", func(ctx *noggo.TriggerContext) {
		ctx.Finish()
	})
	//添加一个任务，noggo.Task.After("task.test", time.Second*3) 3秒后执行任务
	noggo.Add("task.test", func(ctx *noggo.TaskContext) {
		ctx.Finish()
	})
	//添加一个每5秒执行的计划，  不需调用，
	noggo.Add("*/5 * * * * *", func(ctx *noggo.PlanContext) {
		noggo.Logger.Debug("定时计划开始执行")
		ctx.Finish()
	})
	//添加一个事件，noggo.Event.Publish("event.test")调用
	noggo.Add("event.test", func(ctx *noggo.EventContext) {
		ctx.Finish()
	})
	//添加一个队列，noggo.Queue.Publish("queue.test")调用
	noggo.Add("queue.test", func(ctx *noggo.QueueContext) {
		ctx.Finish()
	})

	noggo.Launch(":8080")
}
```



# 项目
完整项目示例，请参考 project 目录或子项目 noggo-project
```
go get github.com/nogio/noggo-project
```

