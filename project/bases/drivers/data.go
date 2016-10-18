package drivers


import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/driver/data-pgsql"
	"fmt"
)


func init() {
	//注册pgsql数据驱动
	noggo.Data.Driver("pgsql", data_pgsql.Driver())
	fmt.Printf("注册驱动了啊")
}
