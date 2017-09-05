package http_form

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"path"
	"crypto/sha1"
	"encoding/xml"
	"io"
	"strings"
	"regexp"
)



//返回中间件
func Middler(uploadPaths ...string) (noggo.HttpFunc) {





	uploadPath := os.TempDir()
	if len(uploadPaths) > 0 && uploadPaths[0] != "" {
		uploadPath = uploadPaths[0]
	}

	return func(ctx *noggo.HttpContext) {


		//判断下是否是ajax请求
		//X-Requested-With:XMLHttpRequest
		if ctx.Header("X-Requested-With") == "XMLHttpRequest" {
			ctx.Ajax = true
		}



		//处理表单
		//if ctx.Method == "post" || ctx.Method == "put" || ctx.Method == "patch" || ctx.Method == "delete" {
		if ctx.Method != "get" {

			//根据不同的类型来处理
			contentType := ctx.Req.Header.Get("Content-Type")


			//if contentType == "text/json" || contentType == "application/json" {
			if strings.Index(contentType, "json") >= 0 {
				body, err := ioutil.ReadAll(ctx.Req.Body)
				if err == nil {
					ctx.Raw = string(body)

					m := Map{}
					err := json.Unmarshal(body, &m)
					if err != nil {
						//遍历JSON对象
						for k, v := range m {
							ctx.Form[k] = v
							ctx.Value[k] = v
						}
					}
				}
				//} else if contentType == "text/xml" || contentType == "application/xml" {
			} else if strings.Index(contentType, "xml") >= 0 {
				body, err := ioutil.ReadAll(ctx.Req.Body)
				if err == nil {
					ctx.Raw = string(body)

					m := Map{}
					err := xml.Unmarshal(body, &m)
					if err != nil {
						//遍历XML对象
						for k,v := range m {
							ctx.Form[k] = v
							ctx.Value[k] = v
						}
					}
				}
			} else {
				ctx.Req.ParseMultipartForm(32 << 20)
			}


			//表单的处理
			//value要考虑    name[a][b][c]
			//这样的话, 就生成了 json值到 ctx.Value 中, 方便一些
			//json[a] json[b] 这样生成: { a:xxx, b: yyy }
			//楼上如果多个   生成: [{a: a1, b: b2}, {a: a2, b: b2}]
			//以上都是计划, 先不管


			for k,v := range ctx.Req.Form {
				//ctx.Form[k] = strings.Join(v, "")
				//ctx.Value[k] = strings.Join(v, "")

				//一个存字串,多个存数组
				if len(v) > 1 {
					ctx.Form[k] = v
					ctx.Value[k] = v
				} else {
					ctx.Form[k] = v[0]
					ctx.Value[k] = v[0]
				}
			}
			for k,v := range ctx.Req.PostForm {
				//ctx.Form[k] = strings.Join(v, "")
				//ctx.Value[k] = strings.Join(v, "")

				//一个存字串,多个存数组
				if len(v) > 1 {
					ctx.Form[k] = v
					ctx.Value[k] = v
				} else {
					ctx.Form[k] = v[0]
					ctx.Value[k] = v[0]
				}
			}

			//处理Post 和 file
			if ctx.Req.MultipartForm != nil {

				//Post值
				for k,v := range ctx.Req.MultipartForm.Value {
					//vv := strings.Join(v, "")
					//ctx.Form[k] = vv
					//ctx.Value[k] = vv


					//有个问题，当type=file时候，又不选文件的时候，value里会存在一个空字串的value
					//如果同一个form name 有多条记录，这时候会变成一个[]string，的空串数组
					//这时候，mapping解析文件的时候[file]就会出问题，会判断文件类型，这时候是[]string就出问题了


					//一个存字串,多个存数组
					if len(v) > 1 {
						ctx.Form[k] = v
						ctx.Value[k] = v
					} else {
						ctx.Form[k] = v[0]
						ctx.Value[k] = v[0]
					}
				}

				//临时保存目录, 若不设置等于系统临时目录
				tempPath := os.TempDir()
				if uploadPath != "" {
					tempPath = uploadPath
				}
				//去掉斜杠
				if tempPath[len(tempPath)-1:] == "/" {
					tempPath = tempPath[0:len(tempPath)-1]
				}


				//FILE可能要弄成JSON，文件保存后，MIME相关的东西，都要自己处理一下
				for k,v := range ctx.Req.MultipartForm.File {

					//这里应该保存为数组
					files := []Map{}

					//处理多个文件
					for _,f := range v {

						hash := ""
						filename := f.Filename
						mimetype := f.Header.Get("Content-Type")
						extension := path.Ext(filename)
						if extension != "" {
							extension = extension[1:]	//去掉.
						}
						var tempfile string
						var length int64

						//先计算hash
						if file, err := f.Open(); err == nil {
							h := sha1.New()
							if _, err := io.Copy(h, file); err == nil {
								hash = fmt.Sprintf("%x", h.Sum(nil))
							}
							file.Close()
						}
						//如果HASH没算出来
						if hash == "" {
							continue
						}


						//保存临时文件
						tempfile = fmt.Sprintf("%s/%s.%s", tempPath, hash, extension)
						if file, err := f.Open(); err == nil {
							if save, err := os.OpenFile(tempfile, os.O_WRONLY|os.O_CREATE, 0777); err == nil {
								io.Copy(save, file)	//保存文件
								save.Close()
							}
							file.Close()
						}

						//读文件大小信息
						if fi, err := os.Stat(tempfile); err == nil {
							length = fi.Size()
						}

						//fmt.Printf("k=%s, hash=%v, name=%v, ext=%v, mime=%v, length=%v, tempfile=%v\n", k, hash, filename, extension, mimetype, length, tempfile)

						if length == 0 {
							continue
						}

						msg := Map{
							"hash": hash,
							"filename": filename,
							"extension": strings.ToLower(extension),
							"mimetype": mimetype,
							"length": length,
							"tempfile": tempfile,
						}

						files = append(files, msg)

						/*
						ctx.Value[k] = msg
						ctx.Files[k] = msg
						*/


						//fmt.Printf("k=%v,v=%v\n", k, msg)


						/*
	
						if file, err := f.Open(); err == nil {
							defer file.Close()
							h := sha1.New()
							if _, err := io.Copy(h, file); err == nil {
	
								hash = fmt.Sprintf("%x", h.Sum(nil))
								fmt.Printf("hs=%v\n", h.Size())
	
								//保存文件
								tempfile := fmt.Sprintf("%s/%s.%s", tempPath, hash, extension)
								if save, err := os.OpenFile(tempfile, os.O_WRONLY|os.O_CREATE, 0777); err == nil {
									//io.Copy(save, file)	//保存文件, file已经被sha1时读走了
									//io.Copy(save, h)	//保存文件
									save.Close()
								} else {
									//保存文件失败,跳过吧
									continue
								}
	
								//读取文件大小信息
								/*
								fi, _ := os.Stat(tempfile)
								if fi != nil && !fi.IsDir() {
									length = fi.Size()
								} else {
									//文件信息拿不到,跳过
									continue
								}
	
								fmt.Printf("k=%s, hash=%v, name=%v, ext=%v, mime=%v, length=%v, tempfile=%v\n", k, hash, filename, extension, mimetype, length, tempfile)
							}
						}
						*/
					}

					//单个单个。 多个数组
					if len(files) > 1 {
						ctx.Upload[k] = files
						ctx.Value[k] = files
					} else {
						ctx.Upload[k] = files[0]
						ctx.Value[k] = files[0]
					}

				}

			}
		}

		//对value进行json的处理，如下
		//user[id] user[name] 解析为     user: { id, name }
		//多个以上的，则解析为 user: [{id,name},{id,name}]

		jsons := map[string]Map{}
		lists := map[string][]Map{}


		for k,v := range ctx.Value {
			regx := regexp.MustCompile(`^([a-zA-Z0-9_-]+)\[([a-zA-Z0-9_-]+)\]$`)
			if regx.MatchString(k) {

				matchs := regx.FindAllStringSubmatch(k, -1)
				if matchs != nil && len(matchs) > 0 && len(matchs[0]) > 0 {

					obj,key := matchs[0][1],matchs[0][2]

					//如果值是string，就是单个实体
					//如果值是[]string，就是实体数组
					switch sv := v.(type) {
					case string:
						if jsons[obj] != nil {
							jsons[obj][key] = sv
						} else {
							jsons[obj] = Map{
								key: sv,
							}
						}
					case []string:
						if lists[obj] != nil {
							maps := lists[obj]
							for svIndex,svVal := range sv {
								maps[svIndex][key] = svVal
							}
						} else {
							maps := []Map{}
							for _,svVal := range sv {
								maps = append(maps,Map{
									key: svVal,
								})
							}
							lists[obj] = maps
						}
					}
				}
			}
		}

		for k,v := range jsons {
			ctx.Value[k] = v
		}
		for k,v := range lists {
			ctx.Value[k] = v
		}

		ctx.Next()

	}
}
