package data_postgres_relate

/*
	此驱动为内部使用，所有表/模型(model），必须有以下3个字段：
	count   jsonb       记录父表的统计数量，如：{"logins": 0}
	status  varchar     状态，为null表示正常，非null表示逻辑删除
	changed datetime    最后修改时间
	created datetime    创建时间

	定义模型的时候，字段使用，relate 来配置父子关系，如当前为login表，父表为user表：
	"user_id": Map{
		"type": "int", "must": true, "name": "所属用户编号", "text": "所属用户编号",
		"relate": Map{ "parent": "user", "count": "logins", "status": "user.removed" }
	}
	以上定义表示，
	parent指向父表model，
	count表示父表count字段下的logins子字段，当子表添加记录时，父表关联记录的count字段中的logins子字段会被+1，以冗余子表记录数
	status表示，当父表user记录被删除时
	子表的相关记录的status会被更新为user.removed，表示逻辑删除，当父表记录被恢复时，子表status为user.remove的记录同步被恢复

*/


import (
	_ "github.com/nogio/noggo/depend/pq"   //此包自动注册名为postgres的sql驱动
	"github.com/nogio/noggo"
)


const (
	SQLDRIVER = "postgres"

	FieldChanged = "changed"

	StatusField = "status"
	StatusRemoved = "removed"

	TriggerCreate = "data.create"
	TriggerChange = "data.change"
	TriggerRemove = "data.remove"
	TriggerRecover = "data.recover"
)


//返回驱动
func Driver() (noggo.DataDriver) {
	return &PostgresDriver{}
}