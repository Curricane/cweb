# 介绍
本项目是类似gin一样但更小的http框架，支持trie树实现的动态路由，上下文，中间件等基础功能。
## 目录结构
```
cweb
├── context.go // 上下文
├── cweb.go // 框架入口
├── go.mod 
├── logger.go // 日志
├── readme.md
├── recovery.go // 错误恢复
├── router.go // 路由
└── trie.go // trie树，用于实现动态路由
```
## 使用例子
```golang
package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/Curricane/cweb"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	//r := cweb.Default() // 使用默认的Engine
	r := cweb.New()      // 使用无中间件的Engine
	r.Use(cweb.Logger()) // 添加Logger中间件
	r.SetFuncMap(template.FuncMap{ // SetFuncMap for custom render function
		"FormatAsDate": FormatAsDate,
	})

	r.LoadHTMLGlob("templates/*")   // 加载HTMLGlob
	r.Static("/assets", "./static") // 静态文件

	stu1 := &student{Name: "cc", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}

	r.GET("/", func(c *cweb.Context) { //添加路由
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/students", func(c *cweb.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", cweb.H{
			"title":  "cweb",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/date", func(c *cweb.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", cweb.H{
			"title": "cweb",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	r.Run(":9999") // 运行
}
```
## 实现细节
- Context 封装了http 的请求和响应，并提供了便利的接口，如JSON，Date
- 分组路由的实现由Engine开始就内嵌了RouterGroup
- 动态路由由trie树来实现，可以快速的查找，并处理: *等模式的动态路由
