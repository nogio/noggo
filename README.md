# 注意

此框架刚完成，尚未正式发布，目前我们正在一些正式项目中使用此框架，以便对框架进行细节的打磨和完善，请暂时不要在生产环境使用此框架，可能还会有很多调整和优化的地方。

# 简介

noggo 是 隶属于 nogio 基于 golang 实现的后端开发框架，如框架名字(nog(木栓子、木楔子))一样，我们希望成为后端开发的连接件，整合各大模块，让开发者可以专注于业务开发，而非浪费大量时间在基本件上的处理。

# 目标

noggo所追求的是开发效率，而单节点的极致性能不是noggo所追求的，若要提升性能可以多节点负载均衡等多种手段实现。
noggo整合后端开发常用的，日志、触发器、任务、计划、订阅/发布、队列，数据库（SQL,NOSQL)、缓存等模块，并提供一致性的开发体验，降低后端开发门槛，提升后端开发效率。
noggo可以让您使用一套框架，完成整个后端的开发工作，并提供自动生成文档等功能，为后端开发与前端开发的对接提供无限便利。


# 联络

官方QQ群：34613300

# 特色

*   驱动化
*   中间件
*   多语言



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
	_ "github.com/nogio/noggo/core" //引用框架默认驱动等等，也可自己定义各种驱动
	"github.com/nogio/noggo/middler"
)



func main() {

	//请求日志与静态文件中间件
	noggo.Use(middler.HttpLogger())
	noggo.Use(middler.HttpStatic("statics"))

	//get 首页
	noggo.Get("/", func(ctx *noggo.HttpContext) {
		ctx.Text("hello noggo.")
	})
	//post 首页
	noggo.Post("/", func(ctx *noggo.HttpContext) {
		ctx.Text("post hello noggo.")
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
完整项目示例，请参考 noggo-project 
```
go get github.com/nogio/noggo-project
```

