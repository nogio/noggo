package main

import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/middler"
)

func main() {

	nog := noggo.New()

	//请求日志与静态文件中间件
	nog.Use(middler.HttpLogger())
	nog.Use(middler.HttpStatic("statics"))

	//Get请求首页
	nog.Get("/", func(ctx *noggo.HttpContext) {
		ctx.Data["msg"] = "这是什么消息"
		ctx.View("index")

		//3秒后开始一个任务
		//noggo.Task.After("test", time.Second*3)

		//返回一段文本给客户端
		//ctx.Text("hello noggo")
	})


	//添加一个每10秒运行的周期性计划
	nog.Add("*/10 * * * * *", func(ctx *noggo.PlanContext) {
		noggo.Logger.Debug("10秒计划开始执行了")
		ctx.Finish()
	})

	//添加一个测试任务
	nog.Add("test", func(ctx *noggo.TaskContext) {
		noggo.Logger.Debug("测试任务开始执行了")
		ctx.Finish()
	})



	nog.Run(":8080")
}