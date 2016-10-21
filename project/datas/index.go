package data


import (
	_ "./bases"
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/driver/data-pgsql"
)

func init() {
	noggo.Current = ""

	//注册数据层驱动
	noggo.Data.Driver("pgsql", data_pgsql.Driver())
}