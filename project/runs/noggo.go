package main

import (
	"github.com/nogio/noggo"
	_ "github.com/nogio/noggo/core"
	_ "../bases"
	_ "../datas"
	_ "../globals"
	_ "../modules"
	_ "../nodes"
	"os"
	"github.com/nogio/noggo/driver/data-postgres"
	"github.com/nogio/noggo/driver/data-mysql"
	"github.com/nogio/noggo/driver/data-adodb"
	"github.com/nogio/noggo/driver/data-sqlite"
)

func init() {
	//基础驱动和默认方法在  github.com/nogio/noggo/core 包中
	//直接引用即可， 否则所有驱动，以及类型，加密方法等等都需要手动注册
	noggo.Driver("postgres", data_postgres.Driver())
	noggo.Driver("mysql", data_mysql.Driver())
	noggo.Driver("adodb", data_adodb.Driver())
	noggo.Driver("sqlite", data_sqlite.Driver())
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
