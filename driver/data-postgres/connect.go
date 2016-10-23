package data_postgres

import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"database/sql"
	"errors"
)

type (
	//数据库连接
	PostgresConnect struct {
		config Map

		//数据库对象
		db  *sql.DB
		//连接串
		url string
		//模型
		models  map[string]Map
	}
)

//打开连接
func (conn *PostgresConnect) Open() error {
	db, err := sql.Open(SQLDRIVER, conn.url)
	if err != nil {
		return errors.New("数据库连接失败：" + err.Error())
	} else {
		conn.db = db
		return nil
	}
}
//关闭连接
func (conn *PostgresConnect) Close() error {
	if conn.db != nil {
		err := conn.db.Close()
		conn.db = nil
		return err
	}
	return nil
}


//注册模型
func (conn *PostgresConnect) Model(name string, config Map) {
	conn.models[name] = config
}

func (conn *PostgresConnect) Base(name string, cache driver.CacheBase) (driver.DataBase,error) {
	return &PostgresBase{name, conn, conn.models, conn.db, nil, cache, false},nil
}