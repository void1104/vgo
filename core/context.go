package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
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
	Handlers []HandlerFunc
	index    int

	// This mutex protect Keys map
	mu sync.RWMutex
	// keys is a k/v pair exclusively for the context of each request.
	keys map[string]interface{}
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

// Fail 返回失败状态
func (c *Context) Fail() {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(500)
	c.Writer.Write([]byte("Internal Server Error"))
}

// AuthFail 返回鉴权失败状态
func (c *Context) AuthFail() {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(403)
	c.Writer.Write([]byte("Forbidden, Auth Fail"))
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// Set is used to store a k/v pair exclusively for this context.
// It also lazy initializes c.Keys if it was not used previously.
func (c *Context) Set(key string, value interface{}) {
	c.mu.Lock()
	if c.keys == nil {
		c.keys = make(map[string]interface{})
	}

	c.keys[key] = value
	c.mu.Unlock()
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (c *Context) Get(key string) (value interface{}, exists bool) {
	c.mu.RLocker()
	value, exists = c.keys[key]
	c.mu.RUnlock()
	return
}

// GetString returns the value associated with the key as a string.
func (c *Context) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

// GetBool returns the value associated with the key as a boolean.
func (c *Context) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

// GetInt returns the value associated with the key as an integer.
func (c *Context) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

// GetInt64 returns the value associated with the key as an integer.
func (c *Context) GetInt64(key string) (i int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int64)
	}
	return
}

// GetUint returns the value associated with the key as an unsigned integer.
func (c *Context) GetUint(key string) (i uint) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(uint)
	}
	return
}

// GetUint64 returns the value associated with the key as an unsigned integer.
func (c *Context) GetUint64(key string) (i uint64) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(uint64)
	}
	return
}

// GetFloat64 returns the value associated with the key as an float64.
func (c *Context) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

// GetTime returns the value associated with the key as a time.
func (c *Context) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}
