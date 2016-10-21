package logics

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
)






func GetTest(id Any) (error,Map) {
	db := noggo.Data.Base("main"); defer db.Close()
	return db.Model("test").Entity(id)
}


func AddTest(data Map) (error,Map) {
	db := noggo.Data.Base("main"); defer db.Close()
	return db.Model("test").Create(data)
}


func ChangeTest(item, data Map) (error,Map) {
	db := noggo.Data.Base("main"); defer db.Close()
	return db.Model("test").Change(item, data)
}


func RemoveTest(item Map) (error) {
	db := noggo.Data.Base("main"); defer db.Close()
	return db.Model("test").Remove(item)
}

