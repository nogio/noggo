package data_mysql

import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"database/sql"
	"errors"
)

type (
	//数据库连接
	MysqlConnect struct {
		config Map

		//数据库对象
		db  *sql.DB
		//连接串
		url string
		//模型
		models  map[string]Map
		views  map[string]Map
	}
)

//打开连接
func (conn *MysqlConnect) Open() error {
	db, err := sql.Open(SQLDRIVER, conn.url)
	if err != nil {
		return errors.New("数据库连接失败：" + err.Error())
	} else {
		conn.db = db
		return nil
	}
}
//关闭连接
func (conn *MysqlConnect) Close() error {
	if conn.db != nil {
		err := conn.db.Close()
		conn.db = nil
		return err
	}
	return nil
}


//注册模型
func (conn *MysqlConnect) Model(name string, config Map) {
	conn.models[name] = config
}
func (conn *MysqlConnect) View(name string, config Map) {
	conn.views[name] = config
}

func (conn *MysqlConnect) Base(name string, cache noggo.CacheBase) (noggo.DataBase,error) {
	return &MysqlBase{name, conn, conn.models, conn.views, conn.db, nil, cache, true, false},nil
}




//来来来，构建数据结构
func (conn *MysqlConnect) Build() error {
	return nil
}