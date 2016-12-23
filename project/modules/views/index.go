package helpers

/*
    此包注册在view视图层使用的函数
*/

import (
    "github.com/nogio/noggo"
    . "github.com/nogio/noggo/base"
    "fmt"
    "strings"
    "html/template"
)

func init() {

    noggo.View.Helper("test", func() string {
        return "test"
    })

    //显示文件大小
    noggo.View.Helper("length", func(length int64) string {

        l := float64(length)

        if l > (1024*1024*1024) {
            return fmt.Sprintf("%.2f GB", (l/1024/1024/1024))
        } else if l > (1024*1024) {
            return fmt.Sprintf("%.2f MB", (l/1024/1024))
        } else if l > (1024) {
            return fmt.Sprintf("%.2f KB", (l/1024))
        } else {
            return fmt.Sprintf("%.2f Bytes", l)
        }
    })


    //高亮显示
    noggo.View.Helper("light", func(t, k Any) template.HTML {

        text := fmt.Sprintf("%v", t)
        keyword := fmt.Sprintf("%v", k)

        if keyword == "" {
            return template.HTML(text)
        }
        html := strings.Replace(text, keyword, fmt.Sprintf(`<span style="color:red;">%v</span>`, keyword), -1)
        return template.HTML(html)
    })



}
