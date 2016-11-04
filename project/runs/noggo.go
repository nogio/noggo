package main

import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/driver"
	"github.com/nogio/noggo/middler"
	_ "github.com/nogio/noggo/core"
	_ "../bases"
	_ "../datas"
	_ "../modules"
	_ "../nodes"
	"os"
)

func init() {
	//驱动
	noggo.Driver("postgres", driver.DataPostgres())
	noggo.Driver("mysql", driver.DataMysql())
	noggo.Driver("adodb", driver.DataAdodb())
	noggo.Driver("sqlite", driver.DataSqlite())

	//中间件
	noggo.Use(middler.HttpLogger())
	noggo.Use(middler.HttpStatic("statics"))
	noggo.Use(middler.HttpForm("uploads"))
}

func main() {
	//框架初始化
	noggo.Init()


	//可以用命令行运行多个节点在同一进程
	//noggo				运行所有节点
	//noggo	*			运行所有节点
	//noggo www ing		运行指定节点


	//同时运行多个/所有节点
	nodes := []string{}
	if len(os.Args) == 1 || (len(os.Args) == 2 && os.Args[1] == "*") {
		for k, _ := range noggo.Config.Node {
			nodes = append(nodes, k)
		}
	} else {
		for i, s := range os.Args {
			if i > 0 {
				nodes = append(nodes, s)
			}
		}
	}
	//开始运行
	for _, name := range nodes {
		node := noggo.New(name)
		if node != nil {
			node.Run()
		}
	}


	//监听退出
	noggo.Exit()
}
