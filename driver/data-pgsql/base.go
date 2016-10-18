package data_pgsql


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"database/sql"
	"errors"
)

type (
	PgsqlBase struct {
		name    string
		conn    *PgsqlConnect
		models  map[string]Map

		db      *sql.DB
		tx      *sql.Tx
	}
)


//关闭数据库
func (base *PgsqlBase) Close() {

	//好像目前不需要关闭什么东西
	if base.tx != nil {
		//关闭时候,一定要提交一次事务
		//如果手动提交了, 这里会失败, 问题不大
		//如果没有提交的话, 连接不会交回连接池. 会一直占用
		base.Cancel()
	}
}

//获取模型对象
func (base *PgsqlBase) Model(name string) (driver.DataModel) {
	if config,ok := base.models[name]; ok {

		//模式，表名
		schema, table, key, fields := "public", name, "id", Map{}
		if n,ok := config["schema"].(string); ok {
			schema = n
		}
		if n,ok := config["table"].(string); ok {
			table = n
		}
		if n,ok := config["key"].(string); ok {
			key = n
		}
		if n,ok := config["fields"].(Map); ok {
			fields = n
		}

		return &PgsqlModel{ base, name, schema, table, key, fields }
	} else {
		panic("数据：模型不存在")
	}
}



//开启事务
func (base *PgsqlBase) Begin() (error, *sql.Tx) {

	if base.tx == nil {
		tx,err := base.db.Begin()
		if err != nil {
			return err, nil
		}
		base.tx = tx
	}
	return nil, base.tx
}

//提交事务
func (base *PgsqlBase) Submit() (error) {

	if base.tx == nil {
		return errors.New("数据：tx未开始")
	}

	err := base.tx.Commit()
	if err != nil {
		return err
	}

	//提交后,要清掉事务
	base.tx = nil

	//提交事务后,要把触发器都发掉
	/*
	for _,trigger := range base.triggers {
		Trigger.Touch(trigger.name, trigger.value)
	}
	*/

	return nil
}

//取消事务
func (base *PgsqlBase) Cancel() (error) {

	if base.tx == nil {
		return errors.New("数据：tx未开始")
	}

	err := base.tx.Rollback()
	if err != nil {
		return err
	}

	//取消后,要清掉事务
	base.tx = nil

	return nil
}