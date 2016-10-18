package data_pgsql


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"errors"
)

type (
	PgsqlDriver struct {

	}
)

//驱动连接
func (drv *PgsqlDriver) Connect(config Map) (error,driver.DataConnect) {

	if config == nil {
		return errors.New("配置不可为空"),nil
	}

	if url,ok := config["url"].(string); ok {

		return nil,&PgsqlConnect{
			config: config, url: url, models: map[string]Map{},
		}

	} else {
		return errors.New("配置缺少[url]信息"),nil
	}
}
