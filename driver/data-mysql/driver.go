package data_mysql


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"errors"
	"net/url"
	"strings"
)

type (
	MysqlDriver struct {

	}
)

//驱动连接
func (drv *MysqlDriver) Connect(config Map) (noggo.DataConnect,error) {

	if config == nil {
		return nil,errors.New("配置不可为空")
	}

	if urlstr,ok := config["url"].(string); ok {

		schema := ""
		//解析url，得到真正的库名，用于schema
		//因为mysql的schema就是库名，灰机的~~~~~


		obj,err := url.Parse(urlstr)
		if err == nil {
			if obj.Path != "" {
				//因为path是从 / 开始的
				schema = obj.Path[0:]
			} else {
				schema = urlstr[strings.Index(urlstr, "/")+1:strings.Index(urlstr, "?")]
			}
		}

		return &MysqlConnect{
			config: config, url: urlstr, schema: schema, models: map[string]Map{}, views: map[string]Map{},
		},nil

	} else {
		return nil,errors.New("配置缺少[url]信息")
	}
}
