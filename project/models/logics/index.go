package logics

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
    "github.com/nogio/noggo/project/mains/consts"
    "time"
)



//按帐号查询管理员
func GetAdminByAccount(account string) (Map,error) {
    db := noggo.Data.Base(consts.MAINDB); defer db.Close()
    return db.Model("admin").Single(Map{
        "account": account,
    })
}

//查询所有管理员
func GetAdmins() ([]Map,error) {
    db := noggo.Data.Base(consts.MAINDB); defer db.Close()
    return db.Model("admin").Query()
}


//添加管理员
func NewAdmin(data Map) (Map,error) {
    db := noggo.Data.Base(consts.MAINDB); defer db.Close()
    return db.Model("admin").Create(data)
}


//修改管理员信息
func ChangeAdmin(item, data Map) (Map,error) {
    db := noggo.Data.Base(consts.MAINDB); defer db.Close()
    data["changed"] = time.Now()
    return db.Model("admin").Change(item, data)
}


//删除管理员
func RemoveAdmin(item Map) (error) {
    db := noggo.Data.Base(consts.MAINDB); defer db.Close()
    return db.Model("admin").Remove(item)
}