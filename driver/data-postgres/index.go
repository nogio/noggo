package data_postgres

import (
	_ "github.com/nogio/noggo/depend/pq"   //此包自动注册名为postgres的sql驱动
	"github.com/nogio/noggo"
)


const (
	SQLDRIVER = "postgres"
)


//返回驱动
func Driver() (noggo.DataDriver) {
	return &PostgresDriver{}
}