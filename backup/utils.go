package noggo

import (
	. "github.com/nogio/noggo/base"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

func readJsonFile(filename string) (Map) {
	bytes, err := ioutil.ReadFile(filename)
	if err == nil {
		m := make(Map)
		err := json.Unmarshal(bytes, &m)
		if err == nil {
			return m
		}
	}
	return nil
}



func readJsonConfig() (configConfig) {
	m := configConfig{}

	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(fmt.Sprintf("读取配置文件出错:%v", err))
	} else {
		err := json.Unmarshal(bytes, &m)
		if err != nil {
			panic(fmt.Sprintf("解析配置文件出错:%v", err))
		}
	}

	return m
}
