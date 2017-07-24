package data_cockroach


import (
	. "github.com/nogio/noggo/base"
	"errors"
)

type (
	CockroachModel struct {
		base    *CockroachBase
		name    string  //模型名称
		schema  string  //架构名
		object   string  //这里可能是表名，视图名，或是集合名（mongodb)
		key     string  //主键
		fields  Map     //字段定义
	}
)


//查询单条
func (model *CockroachModel) Single(args ...Any) (Map,error) {
	return nil,errors.New("有个问题，golang的map顺序不是固定的，这不好处理啊，正在想办法中")
}
//查询列表
func (model *CockroachModel) Query(args ...Any) ([]Map,error) {
	return nil,errors.New("有个问题，golang的map顺序不是固定的，这不好处理啊，正在想办法中")
}



