/*
	此文件代码只是语法糖
	提供一套简单使用框架的方法
	正式项目不建议如此使用
*/

package main

import (
	"fmt"
	"strings"
	"errors"
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
)

func main() {

	sql,args,err := Parsing(Map{
		"status": nil, "$or1": []Map{
			Map{ "name": noggo.FTS("ha") },
			Map{ "stops": noggo.FTS("ha") },
		},
	}, Map{
		"id": noggo.ASC, "views": noggo.DESC,
	})

	fmt.Println(sql,args,err)

	/*
	r1,r2,r3 := parsing(Map{
		"name": "yang", "status": nil, "intro": noggo.FTS("中国"), "test": noggo.NotNil{},
	},Map{
		"role": "haha", "status": nil, "views": Map{ ">": 10, "<": 20 },
		"$or1": []Map{
			Map{ "name": noggo.FTS("ha"),"text": noggo.FTS("ha") }, Map{ "text": noggo.FTS("ha") },
		},
		"$or2": []Map{
			Map{ "name": noggo.FTS("ha") }, Map{ "text": noggo.FTS("ha") },
		},
	}, Map{
		"id": noggo.ASC, "views": noggo.DESC,
	})

	fmt.Println(r1,r2,r3)
	*/

}



//查询语法解析器
// 所有参数使问号(?)
// 所有字段使用$包裹，示例如下：
// $name$=? AND $role$=?,["yang"],nil
// 写数据驱动时，请自行处理字段包裹和参数
// 如postgresql使用的是 $1,$2，字段使用引号""包裹，需要自己处理
// mysql 字段使用``包裹，需要自行处理等等
// mongodb 自行处理，不需要使用此解析器
func Parsing(args ...Any) (string,[]interface{},error) {

	if len(args) > 0 {

		//如果直接写sql
		if v,ok := args[0].(string); ok {
			sql := v
			params := []interface{}{}

			for i,arg := range args {
				if i > 0 {
					params = append(params, arg)
				}
			}

			return sql,params,nil

		} else {

			maps := []Map{}
			for _,v := range args {
				if m,ok := v.(Map); ok {
					maps = append(maps, m)
				}
			}

			querys,values,orders := parsing(maps...)

			orderStr := ""
			if len(orders) > 0 {
				orderStr = fmt.Sprintf("ORDER BY %s", strings.Join(orders, ","))
			}

			sql := fmt.Sprintf("%s %s", strings.Join(querys, " OR "), orderStr)

			return sql, values, nil


		}
	}

	return "",[]interface{}{},errors.New("解析失败")
}

func parsing(args ...Map) ([]string,[]interface{},[]string) {

	fp := ""

	querys := []string{}
	values := make([]interface{}, 0)
	orders := []string{}

	//否则是多个map,单个为 与, 多个为 或
	for _,m := range args {
		ands := []string{}

			for k,v := range m {

				//如果值是ASC,DESC，表示是排序
				if ov,ok := v.(string); ok && (ov==noggo.ASC || ov==noggo.DESC) {

					if ov == noggo.ASC {
						orders = append(orders, fmt.Sprintf(`%s%s%s ASC`, fp, k, fp))
					} else {
						orders = append(orders, fmt.Sprintf(`%s%s%s DESC`, fp, k, fp))
					}

				} else if ms,ok := v.([]Map); ok {
					//是[]Map，相当于or

					qs,vs,os := parsing(ms...)
					if len(qs) > 0 {
						ands = append(ands, fmt.Sprintf("(%s)", strings.Join(qs, " OR ")))
						for _,vsVal := range vs {
							values = append(values, vsVal)
						}
					}
					for _,osVal := range os {
						orders = append(orders, osVal)
					}




				} else {

					//v要处理一下如果是map要特别处理
					//key做为操作符，比如 > < >= 等
					//而且多个条件是and，比如 views > 1 AND views < 100
					//自定义操作符的时候，可以用  is not null 吗？
					if opMap, opOK := v.(Map); opOK {

						opAnds := []string{}
						for opKey,opVal := range opMap {
							opAnds = append(opAnds, fmt.Sprintf(`%s%s%s %s ?`, fp, k, fp, opKey))
							values = append(values, opVal)
						}
						ands = append(ands, fmt.Sprintf("(%s)", strings.Join(opAnds, " AND ")))

					} else {

						if v == nil {
							ands = append(ands, fmt.Sprintf(`%s%s%s IS NULL`, fp, k, fp))
						} else if _,ok := v.(noggo.Nil); ok {
							//为空值
							ands = append(ands, fmt.Sprintf(`%s%s%s IS NULL`, fp, k, fp))
						} else if _,ok := v.(noggo.NotNil); ok {
							//不为空值
							ands = append(ands, fmt.Sprintf(`%s%s%s IS NOT NULL`, fp, k, fp))
						} else if fts,ok := v.(noggo.FTS); ok {
							//处理模糊搜索
							safeFts := strings.Replace(string(fts), "'", "''", -1)
							ands = append(ands, fmt.Sprintf(`%s%s%s LIKE '%%%s%%'`, fp, k, fp, safeFts))
						} else {
							ands = append(ands, fmt.Sprintf(`%s%s%s = ?`, fp, k, fp))
							values = append(values, v)
						}
					}

				}

		}

		if len(ands) > 0 {
			querys = append(querys, fmt.Sprintf("(%s)", strings.Join(ands, " AND ")))
		}
	}

	return querys,values,orders
}









/*
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
	noggo.Add("0 0 * * * *", func(ctx *noggo.PlanContext) {
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

*/