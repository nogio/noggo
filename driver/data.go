/*
	数据接口
	2016-10-22  定稿
*/

package driver

import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/driver/data-adodb"
	"github.com/nogio/noggo/driver/data-mysql"
	"github.com/nogio/noggo/driver/data-postgres"
	"github.com/nogio/noggo/driver/data-sqlite"
)



//adobe数据驱动
func DataAdodb() (noggo.DataDriver) {
	return data_adodb.Driver()
}


//Mysql数据驱动
func DataMysql() (noggo.DataDriver) {
	return data_mysql.Driver()
}

//postgres数据驱动
func DataPostgres() (noggo.DataDriver) {
	return data_postgres.Driver()
}


//Sqlite数据驱动
func DataSqlite() (noggo.DataDriver) {
	return data_sqlite.Driver()
}