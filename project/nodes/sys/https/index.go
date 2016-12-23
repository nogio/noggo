package https

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"../../../mains/consts"
	"../../../models/logics"
)

func init() {

	noggo.Http.Route("index", Map{
		"uri": "/",
		"auth": Map{
			"admin": Map{
				"sign": "admin", "must": true, "base": consts.MAINDB, "model": "admin", "name": "管理员", "text": "管理员记录",
			},
		},
		"route": Map{
			"name": "后台首页", "text": "后台首页",
			"action": func(ctx *noggo.HttpContext) {
				ctx.View("index/index")
			},
		},
	})



	noggo.Http.Route("login", Map{
		"uri": "/login",
		"route": Map{
			"get": Map{
				"name": "后台登录", "text": "后台登录",
				"action": func(ctx *noggo.HttpContext) {
					ctx.View("index/login")
				},
			},
			"post": Map{
				"name": "后台登录提交", "text": "后台登录提交",
				"args": Map{
					"account": Map{
						"type": "string", "must": true, "name": "帐号", "text": "帐号",
					},
					"password": Map{
						"type": "password", "must": true, "name": "密码", "text": "密码",
					},
				},
				"action": func(ctx *noggo.HttpContext) {
					item,_ := logics.GetAdminByAccount(ctx.Args["account"].(string))
					if item == nil {
						ctx.Alert("帐号不存在")
					} else {

						if item["password"] != ctx.Args["password"] {
							ctx.Alert("密码错误")
						} else {
							//登录成功,签入信息
							ctx.Sign.In("admin", item["id"], item["name"])
							ctx.Goback()
						}
					}
				},
			},
		},
	})


	noggo.Http.Route("logout", Map{
		"uri": "/logout",
		"auth": noggo.Defined("sys.auth"),
		"route": Map{
			"name": "后台登出", "text": "后台登出", "coded": true,
			"action": func(ctx *noggo.HttpContext) {
				//签出
				ctx.Sign.Out("admin")
				ctx.Goback()
			},
		},
	})



	noggo.Http.Route("password", Map{
		"uri": "/password",
		"auth": Map{
			"admin": Map{
				"sign": "admin", "must": true, "base": consts.MAINDB, "model": "admin", "name": "管理员", "text": "管理员记录",
			},
		},
		"route": Map{
			"get": Map{
				"name": "修改密码", "text": "修改密码",
				"action": func(ctx *noggo.HttpContext) {
					ctx.View("index/password")
				},
			},
			"post": Map{
				"name": "修改密码", "text": "修改密码", "code": true,
				"args": Map{
					"current": Map{
						"type": "password", "must": false,  "name": "当前密码", "text": "当前密码,可为空,第一次设置不需要当前密码",
					},
					"password": Map{
						"type": "password", "must": true, "name": "新密码", "text": "新密码",
					},
					"confirm": Map{
						"type": "password", "must": true, "name": "确认密码", "text": "确认密码",
					},
				},
				"action": func(ctx *noggo.HttpContext) {
					item := ctx.Auth["admin"].(Map)

					if ctx.Args["password"] != ctx.Args["confirm"] {
						ctx.Alert("两次密码输入不一致。")
					} else {


						if item["password"] != nil && ctx.Args["current"] != nil && item["password"] != ctx.Args["current"] {
							//当前密码不对
							ctx.Alert("当前密码校验失败。")
						} else {

							_,err := logics.ChangeAdmin(item, Map{
								"password": ctx.Args["password"],
							})
							if err != nil {
								//失败
								ctx.Alert("修改失败")
							} else {

								//成功
								ctx.Show("ok", "修改成功", ctx.Url.Route("index"))

							}
						}

					}


				},
			},
		},
	})

}