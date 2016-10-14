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
		ctx.Text("get hello noggo")
	})

	nog.Run(":8080")
}
```

# 项目
完整项目示例，请参考 nogio/noggo/project 目录

