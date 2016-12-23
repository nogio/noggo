package https

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"../../../models/logics"
	"../../../mains/consts"
)

func init() {


	noggo.Http.Route("admin", Map{
		"uri": "/admin",
		"auth": Map{
			"admin": Map{
				"sign": "admin", "must": true, "base": consts.MAINDB, "model": "admin", "name": "管理员", "text": "管理员登录",
			},
		},
		"route": Map{
			"name": "管理员", "text": "管理员",
			"action": func(ctx *noggo.HttpContext) {
				ctx.Route("admin.filter")
			},
		},
	})

	noggo.Http.Route("admin.filter", Map{
		"uri": "/admin/filter",
		"auth": Map{
			"admin": Map{
				"sign": "admin", "must": true, "base": consts.MAINDB, "model": "admin", "name": "管理员", "text": "管理员登录",
			},
		},
		"route": Map{
			"name": "管理员", "text": "管理员",
			"action": func(ctx *noggo.HttpContext) {

				admins,_ := logics.GetAdmins()

				ctx.Data["admins"] = admins
				ctx.View("admin/filter")
			},
		},
	})


	noggo.Http.Route("admin.create", Map{
		"uri": "/admin/create",
		"auth": Map{
			"admin": Map{
				"sign": "admin", "must": true, "base": consts.MAINDB, "model": "admin", "name": "管理员", "text": "管理员登录",
			},
		},
		"route": Map{
			"get": Map{
				"name": "管理员创建", "text": "管理员创建",
				"action": func(ctx *noggo.HttpContext) {
					ctx.Data["roles"] = noggo.Data.Enums(consts.MAINDB, "admin", "role")
					ctx.View("admin/create")
				},
			},
			"post": Map{
				"name": "管理员创建", "text": "管理员创建",
				"args": noggo.Data.Fields(consts.MAINDB, "admin"),
				"action": func(ctx *noggo.HttpContext) {

					_,err := logics.NewAdmin(ctx.Args)
					if err != nil {
						ctx.Alert("添加失败!")
					} else {
						ctx.Alert("添加成功!", ctx.Url.Back())
					}
				},
			},
		},
	})

	noggo.Http.Route("admin.change", Map{
		"uri": "/admin/change/{id}",
		"auth": Map{
			"admin": Map{
				"sign": "admin", "must": true, "base": consts.MAINDB, "model": "admin", "name": "管理员", "text": "管理员登录",
			},
		},
		"item": Map{
			"admin": Map{
				"args": "id", "base": consts.MAINDB, "model": "admin", "name": "管理员", "text": "管理员记录",
			},
		},
		"args": Map{
			"id": Map{
				"type": "int", "must": true, "name": "编号", "text": "编号",
			},
		},
		"route": Map{
			"get": Map{
				"name": "管理员修改", "text": "管理员修改",
				"action": func(ctx *noggo.HttpContext) {
					item := ctx.Item["admin"].(Map)

					ctx.Data["item"] = item
					ctx.Data["roles"] = noggo.Data.Enums(consts.MAINDB, "admin", "role")
					ctx.View("admin/change")
				},
			},
			"post": Map{
				"name": "管理员修改", "text": "管理员修改",
				"argn": true, "args": noggo.Data.Fields(consts.MAINDB, "admin"),
				"action": func(ctx *noggo.HttpContext) {
					item := ctx.Item["admin"].(Map)

					//密码为空不修改，所以清空ctx.args的password
					if ctx.Args["password"] == nil {
						delete(ctx.Args, "password")
					}

					_,err := logics.ChangeAdmin(item, ctx.Args)
					if err != nil {
						ctx.Alert("修改失败!")
					} else {
						ctx.Alert("修改成功!", ctx.Url.Back())
					}
				},
			},
		},
	})





	noggo.Http.Route("admin.remove", Map{
		"uri": "/admin/remove/{id}",
		"auth": Map{
			"admin": Map{
				"sign": "admin", "must": true, "base": consts.MAINDB, "model": "admin", "name": "管理员", "text": "管理员记录",
			},
		},
		"item": Map{
			"admin": Map{
				"args": "id", "base": consts.MAINDB, "model": "admin", "name": "管理员", "text": "管理员记录",
			},
		},
		"args": Map{
			"id": Map{
				"type": "int", "must": true, "name": "编号", "text": "编号",
			},
		},
		"route": Map{
			"get": Map{
				"name": "管理员删除", "text": "管理员删除",
				"action": func(ctx *noggo.HttpContext) {
					item := ctx.Item["admin"].(Map)

					ctx.Data["item"] = item
					ctx.View("admin/remove")
				},
			},
			"post": Map{
				"name": "管理员删除", "text": "管理员删除",
				"action": func(ctx *noggo.HttpContext) {
					item := ctx.Item["admin"].(Map)

					err := logics.RemoveAdmin(item)
					if err != nil {
						ctx.Alert("删除失败!")
					} else {
						ctx.Alert("删除成功!", ctx.Url.Back())
					}
				},
			},
		},
	})


	noggo.Http.Route("admin.detail", Map{
		"uri": "/admin/detail/{id}",
		"auth": Map{
			"admin": Map{
				"sign": "admin", "must": true, "base": consts.MAINDB, "model": "admin", "name": "管理员", "text": "管理员记录",
			},
		},
		"item": Map{
			"admin": Map{
				"args": "id", "base": consts.MAINDB, "model": "admin", "name": "管理员", "text": "管理员记录",
			},
		},
		"args": Map{
			"id": Map{
				"type": "int", "must": true, "name": "编号", "text": "编号",
			},
		},
		"route": Map{
			"name": "管理员详情", "text": "管理员详情",
			"action": func(ctx *noggo.HttpContext) {
				item := ctx.Item["admin"].(Map)

				ctx.Data["item"] = item
				ctx.View("admin/detail")
			},
		},
	})

}