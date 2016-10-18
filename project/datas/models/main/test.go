package models

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"time"
)

func init() {

	noggo.Data.Register("test", Map{
		"name": "测试", "text": "测试表",
		"fields": Map{
			"id": Map{
				"type": "int", "must": false, "name": "编号", "text": "编号",
			},

			"title": Map{
				"type": "string", "must": true, "name": "标题", "text": "标题",
			},
			"content": Map{
				"type": "string", "must": false, "name": "内容", "text": "内容",
			},

			"changed": Map{
				"type": "datetime", "must": true, "auto": time.Now, "name": "修改时间", "text": "修改时间",
			},
			"created": Map{
				"type": "datetime", "must": true, "auto": time.Now, "name": "创建时间", "text": "创建时间",
			},
		},
	})

}
