/*
	此文件代码只是语法糖
	提供一套简单使用框架的方法
	正式项目不建议如此使用
*/

package main

import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"fmt"
	"strings"
)

func main() {

	/*
	m := Map{
		"index":	ASC,
		"views":	DESC,

		"name": Map{ LIKE: "ASDFASDF" },
	}


	if m["index"] == ASC {
		fmt.Printf("对的")
	} else {
		fmt.Printf("不对")
	}
	*/

	sql,args,order,err := noggo.Data.Parsing(Map{
		"name": NOL, "email": NOTNULL,
		"text": Map{ FULL: "haha", NE: "abcd" },
	}, Map{
		"views": ASC, "id": DESC,
	})

	sql = strings.Replace(sql, DELIMS, `"`, -1)
	order = strings.Replace(order, DELIMS, `"`, -1)

	fmt.Println(err, sql, order, args)

}