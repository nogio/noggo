package maindb

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"time"
)

func init() {

	noggo.Data.Model("admin", Map{
		"name": "管理", "text": "管理表",
		"field": Map{
			"id": Map{
				"type": "int", "must": false, "name": "编号", "text": "编号",
			},

			"account": Map{
				"type": "string", "must": true, "name": "帐号", "text": "帐号",
			},
			"password": Map{
				"type": "password", "must": true, "name": "密码", "text": "密码",
			},
			"name": Map{
				"type": "string", "must": true, "name": "名称", "text": "名称",
			},
			"role": Map{
				"type": "enum", "must": true, "name": "角色", "text": "角色",
				"filter": Map{ "query": true },
				"enum": Map{
					"system":   "管理员",
					"nobody":   "无权限",
				},
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
