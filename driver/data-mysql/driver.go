package data_mysql


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"errors"
)

type (
	MysqlDriver struct {

	}
)

//驱动连接
func (drv *MysqlDriver) Connect(config Map) (driver.DataConnect,error) {

	if config == nil {
		return nil,errors.New("配置不可为空")
	}

	if url,ok := config["url"].(string); ok {
		return &MysqlConnect{
			config: config, url: url, models: map[string]Map{},
		},nil

	} else {
		return nil,errors.New("配置缺少[url]信息")
	}
}
