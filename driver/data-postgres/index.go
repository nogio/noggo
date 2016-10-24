package data_postgres

import (
	"github.com/nogio/noggo/driver"
	_ "github.com/nogio/noggo/driver/data-postgres/pq"   //此包自动注册名为postgres的sql驱动
)


const (
	SQLDRIVER = "postgres"
)


//返回驱动
func Driver() (driver.DataDriver) {
	return &PostgresDriver{}
}