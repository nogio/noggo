package main

import (
	. "github.com/nogio/noggo/base"
	_ "github.com/nogio/noggo/core"
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/middler"
	"github.com/nogio/noggo/driver/data-postgres"
	"time"
)

func init() {
	//注册pgsql数据驱动
	noggo.Data.Driver("postgres", data_postgres.Driver())

	//注册数据模型
	noggo.Data.Model("test", Map{
		"name": "测试", "text": "测试表",
		"schema": "public", "table": "test", "key": "id",   //这行配置可选
		"fields": Map{
			"id": Map{
				"type": "int", "must": false, "name": "编号", "text": "编号",
			},
			"title": Map{
				"type": "string", "must": true, "name": "标题", "text": "标题",
			},
			"content": Map{
				"type": "string", "must": false, "name": "内容", "text": "内容",
			},
			"changed": Map{
				"type": "datetime", "must": true, "auto": time.Now, "name": "修改时间", "text": "修改时间",
			},
			"created": Map{
				"type": "datetime", "must": true, "auto": time.Now, "name": "创建时间", "text": "创建时间",
			},
		},
	})


}


func main() {

	//db := noggo.Data.DB("db")
	//db.Model("user").Query(`abc=$1 and ASDF=$2`, 1, 2)

	/*
	db := noggo.Data.DB("main"); defer db.Close()
	db.Model("user").Entity(1)
	*/



	nog := noggo.New()


	//请求日志与静态文件中间件
	nog.Use(middler.HttpLogger())
	nog.Use(middler.HttpStatic("statics"))

	//Get请求首页
	nog.Get("/", func(ctx *noggo.HttpContext) {


		db := noggo.Data.Base("main"); defer db.Close()

		item,_ := db.Model("test").Create(Map{
			"title": "标题哦", "content": "内容哦",
		})

		item,_ = db.Model("test").Change(item, Map{
			"title": "改过的标题算么", "changed": time.Now(),
		})

		noggo.Event.Publish("test")

		ctx.Json(item)

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


	//添加一个事件
	nog.Add("test", func(ctx *noggo.EventContext) {
		noggo.Logger.Debug("test事件发生")
		ctx.Finish()
	})

	nog.Run(":8080")
}

