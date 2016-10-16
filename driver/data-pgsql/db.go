package data_pgsql


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"database/sql"
)

type (
	PgsqlDB struct {

	}
)


//关闭数据库
func (db *PgsqlDB) Close() {

}

//获取模型对象
func (db *PgsqlDB) Model(name string) (driver.DataModel) {
	return nil
}



//开启事务
func (db *PgsqlDB) Begin() (error, *sql.Tx) {
	return nil,nil
}

//提交事务
func (db *PgsqlDB) Submit() (error) {
	return nil
}

//取消事务
func (db *PgsqlDB) Cancal() (error) {
	return nil
}