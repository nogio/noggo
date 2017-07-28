package data_postgres

import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"database/sql"
	"errors"
)

type (
	//数据库连接
	PostgresConnect struct {
		config Map
		//连接串
		url string
		schema  string

		//数据库对象
		db  *sql.DB

		//模型
		tables  map[string]Map
		views  map[string]Map
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


//注册表
func (conn *PostgresConnect) Table(name string, config Map) {
	conn.tables[name] = config
}

func (conn *PostgresConnect) View(name string, config Map) {
	conn.views[name] = config
}

//注册模型
func (conn *PostgresConnect) Model(name string, config Map) {
	conn.models[name] = config
}

func (conn *PostgresConnect) Base(name string, cache noggo.CacheBase) (noggo.DataBase,error) {
	return &PostgresBase{name, conn, nil, cache, true, false, []noggo.DataTrigger{}},nil
}



//来来来，构建数据结构
func (conn *PostgresConnect) Build() error {
	return nil
}