package data_postgres


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"errors"
)

type (
	PostgresDriver struct {

	}
)

//驱动连接
func (drv *PostgresDriver) Connect(config Map) (noggo.DataConnect,error) {

	if config == nil {
		return nil,errors.New("配置不可为空")
	}

	//支持自定义的schema，相当于数据库名
	schema := SCHEMA
	if vv,ok := config["schema"].(string); ok {
		schema = vv
	}

	if url,ok := config["url"].(string); ok {
		return &PostgresConnect{
			config: config, url: url, schema: schema, db: nil, tables: map[string]Map{}, views: map[string]Map{}, models: map[string]Map{},
		},nil

	} else {
		return nil,errors.New("配置缺少[url]信息")
	}
}
