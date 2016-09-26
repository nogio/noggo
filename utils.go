package noggo

import (
	. "github.com/nogio/noggo/base"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

func readJsonFile(filename string) (m Map) {
	m = make(Map)
	bytes, err := ioutil.ReadFile(filename)
	if err == nil {
		json.Unmarshal(bytes, &m)
	}
	return
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
