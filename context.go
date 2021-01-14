package cweb

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H 取别名
type H map[string]interface{}

// Context 上下文，所有的信息都与Context关联,封装http包，路由，中间件
type Context struct {
	// 封装 http.ResponseWriter *http.Request
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	Params     map[string]string
	StatusCode int

	// 支持中间件
	handlers []HandlerFunc // 中间件
	index    int

	// 支持路由（分组， 动态）
	engine *Engine
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Path:   req.URL.Path,
		Method: req.Method,
		Req:    req,
		Writer: w,
		index:  -1,
	}
}

// Next 调用中间件
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// PostForm 根据键返回表单中的值，若没有返回空string
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 根据键返回url中query中的值
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// String 将按照指定格式字符串传输
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprint(format, values)))
}

// JSON 组成JSON
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	enCoder := json.NewEncoder(c.Writer)
	if err := enCoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}
