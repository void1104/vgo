package context

import (
	"encoding/json"
	"fmt"
	"net/http"
	"vgo/router"
)

/**
对request和response的封装，只是设计Context的原因之一。对于框架来说，还需要支撑额外的功能，
例如保存中间件产生的信息。Context随着每一个请求的出现而产生，请求的结束而销毁，和当前请求强相关的信息都应由
Context承载。因此，设计Context结构，扩展性和复杂性留在了内部，而对外简化了接口。
*/

type H map[string]interface{}

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string
	// response info
	StatusCode int
	// middleware
	Handlers []router.HandlerFunc
	index    int
}

// NewContext context的构造函数
func NewContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (c *Context) Next() {
	c.index++
	s := len(c.Handlers)
	for ; c.index < s; c.index++ {
		c.Handlers[c.index](c)
	}
}

// PostForm  获取请求体中的请求参数 - Form表单
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 获取请求体中的请求参数 - url
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status 设置响应状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader 设置header信息
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// String 返回字符流数据
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON 返回JSON格式数据
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Data 返回字节流数据
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// HTML 返回HTML数据，在目前前后端分离的体系中，不常用
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}
