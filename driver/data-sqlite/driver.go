package data_sqlite


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"errors"
)

type (
	SqliteDriver struct {

	}
)

//驱动连接
func (drv *SqliteDriver) Connect(config Map) (driver.DataConnect,error) {

	if config == nil {
		return nil,errors.New("配置不可为空")
	}

	if file,ok := config["file"].(string); ok {
		return &SqliteConnect{
			config: config, file: file, models: map[string]Map{},
		},nil

	} else {
		return nil,errors.New("配置缺少[file]信息")
	}
}
