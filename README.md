# noggo
noggo Golang 开发框架

# 简介
noggo 是 隶属于 nogio 基于 golang 实现的后端开发框架，如框架名字(nog)一样，我们希望成为后端开发的支架/支撑，让开发者可以专注于业务开发，而非浪费大量时间在基本件上的处理。


# 联络
官方QQ群：34613300


# 示例

```
go get github.com/nogio/noggo
```

```golang
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

		//3秒后开始一个任务
		noggo.Task.After("test", time.Second*3)

        //返回一段文本给客户端
		ctx.Text("hello noggo")
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
```

# 项目
完整项目示例，请参考 nogio/noggo/project 目录

