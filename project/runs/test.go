/*
	此文件代码只是语法糖
	提供一套简单使用框架的方法
	正式项目不建议如此使用
*/

package main

import (
	"github.com/nogio/noggo"
	_ "github.com/nogio/noggo/core" //引用框架默认驱动等等，也可自己定义各种驱动
	"github.com/nogio/noggo/middler"
	"github.com/nogio/noggo/driver/data-postgres"
	"github.com/nogio/noggo/driver/data-mysql"
)

func init() {
	//注册数据层驱动
	noggo.Driver("postgres", data_postgres.Driver())
	noggo.Driver("mysql", data_mysql.Driver())
}


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

