package data_pgsql

import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"database/sql"
	"errors"
)

type (
	//数据库连接
	PgsqlConnect struct {
		config Map

		//数据库对象
		db  *sql.DB
		//连接串
		url string
	}
)

//打开连接
func (conn *PgsqlConnect) Open() error {
	db, err := sql.Open("postgres", conn.url)
	if err != nil {
		return errors.New("数据库连接失败：" + err.Error())
	} else {
		conn.db = db
		return nil
	}
}
//关闭连接
func (conn *PgsqlConnect) Close() error {
	if conn.db != nil {
		err := conn.db.Close()
		conn.db = nil
		return err
	}
	return nil
}


func (conn *PgsqlConnect) DB(name string) (error,driver.DataDB) {
	return nil,nil
}