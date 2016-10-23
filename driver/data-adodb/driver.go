package data_adodb


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"errors"
)

type (
	AdodbDriver struct {

	}
)

//驱动连接
func (drv *AdodbDriver) Connect(config Map) (driver.DataConnect,error) {

	if config == nil {
		return nil,errors.New("配置不可为空")
	}

	if url,ok := config["url"].(string); ok {
		return &AdodbConnect{
			config: config, url: url, models: map[string]Map{},
		},nil

	} else {
		return nil,errors.New("配置缺少[url]信息")
	}
}
