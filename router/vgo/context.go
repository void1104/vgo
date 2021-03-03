package vgo

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/**
目的：
	1. 避免每次请求都去构造完整的request体和response体
	2. 支撑额外的功能，例如：
		- 动态路由/hello/:name中，name的值
		- 中间件产生的信息
		- 请求强相关的信息由context承载
*/

type H map[string]interface{}

type Context struct {
	// 原始对象体
	RW  http.ResponseWriter
	Req *http.Request
	// 请求信息
	Params map[string]string
	// 响应信息
	StatusCode int
}

/**
构造函数
*/
func newContext(rw http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		RW:  rw,
		Req: req,
	}
}

/**
- 在HandlerFunc中，希望能够访问到解析的参数，因此，需要对Context对象增加一个属性和方法，
来提供对路由参数的访问。
- 我们将解析后的参数存储到Params中，通过c.Param("lang")的方式获取到对应的值。
*/
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// 获取Form表单传递字段
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// 获取Get请求字段
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 设置响应状态码
func (c *Context) SetStatus(code int) {
	c.StatusCode = code
	c.RW.WriteHeader(code)
}

// 设置响应头
func (c *Context) SetHeader(key string, value string) {
	c.RW.Header().Set(key, value)
}

// 返回字符串格式数据
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatus(code)
	c.RW.Write([]byte(fmt.Sprintf(format, values...)))
}

// 返回JSON格式数据
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatus(code)
	encoder := json.NewEncoder(c.RW)
	if err := encoder.Encode(obj); err != nil {
		// 当err != nil时，http.Error()不会起作用，因为在WriteHeader()后调用Header().Set()是不会生效的
		//http.Error(c.RW, err.Error(), 500)
		// Gin实现是直接panic
		panic(err)
	}
}

// 返回字节格式数据
func (c *Context) Data(code int, data []byte) {
	c.SetStatus(code)
	c.RW.Write(data)
}

// 返回HTML模板
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatus(code)
	c.RW.Write([]byte(html))
}
