package data_pgsql


import (
	. "github.com/nogio/noggo/base"
	//"github.com/nogio/noggo/driver"
)

type (
	PgsqlModel struct {

	}
)

//创建对象
func (db *PgsqlModel) Create(data Map) (error,Map) {
	return nil, nil
}

//修改对象
func (db *PgsqlModel) Change(item Map, data Map) (error,Map) {
	return nil, nil
}

//删除对象
func (db *PgsqlModel) Remove(Map) (error) {
	return nil
}

//查询对象
func (db *PgsqlModel) Entity(key Any) (error,Map) {
	return nil
}