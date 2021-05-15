package core

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

const abortIndex int8 = math.MaxInt8 / 2

type H map[string]interface{}

// Context is the most important part of gin. It allows us to pass variables
// between middleware, manage the flow, validate the JSON of a request and
// render a JSON response for example
type Context struct {
	//writerMem http.ResponseWriter
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
	index    int8

	// This mutex protect Keys map
	mu sync.RWMutex

	// keys is a k/v pair exclusively for the context of each request.
	keys map[string]interface{}

	// Errors is a list of errors attached to all the handlers/middlewares who used this context.
	Errors errorMsgs

	// Accepted defines a list of manually accepted formats for content negotiation.
	Accepted []string

	// queryCache use url.ParseQuery cached the param query result from c.Request.URL.Query().
	queryCache url.Values

	// formCache use url.ParseQuery cached PostForm contains the parsed form data from POST, PATCH
	// or PUT body parameters.
	formCache url.Values

	// sameSite allows a server to define a cookie attribute making it impossible for
	// the browser to send this cookie along with cross-site requests.
	sameSite http.SameSite
}

/************************************/
/********** CONTEXT CREATION ********/
/************************************/

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

func (c *Context) reset() {
	c.Req = nil
	c.Writer = nil
	c.Path = ""
	c.Method = ""
	c.Params = nil

	c.Handlers = nil
	c.index = -1
	c.keys = nil
}

// Copy returns a copy of the current context that can be safely used outside the request's scope
// This has to be used when the context has to be passed to a goroutine TODO why
func (c *Context) Copy() *Context {
	return nil
}

// HandlerName returns the main handler's name. For example if the handler is 'handlerGetUsers()',
// this function will return 'main.handleGetUses'.
func (c *Context) HandlerName() string {
	return nameOfFunction(c.Handlers[len(c.Handlers)-1])
}

// HandlerNames return a list of all registered handlers for this context in descending order,
// following the semantics of HandlerName()
func (c *Context) HandlerNames() []string {
	hn := make([]string, 0, len(c.Handlers))
	for _, val := range c.Handlers {
		hn = append(hn, nameOfFunction(val))
	}
	return hn
}

// Handler returns the main handler.
func (c *Context) Handler() HandlerFunc {
	return c.Handlers[len(c.Handlers)-1]
}

/************************************/
/*********** FLOW CONTROL ***********/
/************************************/

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.Handlers)) {
		c.Handlers[c.index](c)
		c.index++
	}
}

// IsAborted returns true if the current context was aborted.
func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
// Let's say you have an authorization middleware that validates that the current request is authorized.
// If the authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers
// for this request are not called
func (c *Context) Abort() {
	c.index = abortIndex
}

// AbortWithStatus calls `Abort()` and writes and specified status code.
// For example, a failed attempt to authenticate a request could use: context.AbortWithStatus(401).
func (c *Context) AbortWithStatus(code int) {
	c.Status(code)
	c.Abort()
}

// AbortWithStatusJSON calls `Abort()` and then `JSON` internally.
// This method stops the chain, writes the status code and return a JSON body.
// It also sets the Content-Type as "application/json".
func (c *Context) AbortWithStatusJSON(code int, jsonObj interface{}) {
	c.Abort()
	c.JSON(code, jsonObj)
}

// AbortWithError calls `AbortWithStatus()` and `Error()` internally.
// This method stops the chain, writes the status code and pushes the specified error to `c.Errors`.
// See Context.Error() for more details.
func (c *Context) AbortWithError(code int, err error) *Error {
	c.AbortWithStatus(code)
	return c.Error(err)
}

/************************************/
/********* ERROR MANAGEMENT *********/
/************************************/

// Error attaches an error to the current context. The error is pushed to a list of errors.
// It's a good idea to call Error for each error that occurred during the resolution of a request.
// A middleware can be used to collect all the errors an push them to a database together,
// print a log, or append it in the HTTP response.
// Error will panic if err is nil.
func (c *Context) Error(err error) *Error {
	if err != nil {
		panic("err is nil")
	}

	parsedError, ok := err.(*Error)
	if !ok {
		parsedError = &Error{
			Err:  err,
			Type: ErrorTypePrivate,
		}
	}

	c.Errors = append(c.Errors, parsedError)
	return parsedError
}

/************************************/
/************ INPUT DATA ************/
/************************************/

// PostForm  获取请求体中的请求参数 - Form表单
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 获取请求体中的请求参数 - url
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) initQueryCache() {
	if c.queryCache == nil {
		if c.Req != nil {
			c.queryCache = c.Req.URL.Query()
		} else {
			c.queryCache = url.Values{}
		}
	}
}

func (c *Context) initFormCache() {
	if c.formCache == nil {
		c.formCache = make(url.Values)
		req := c.Req
		c.formCache = req.PostForm
	}
}

// FormFile returns the first file for the provided form key.
func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	if c.Req.MultipartForm == nil {
		return nil, Error{
			Type: ErrorTypePrivate,
			Meta: "request.multipartFrom is nil",
		}
	}

	file, fileHeader, err := c.Req.FormFile(name)
	if err != nil {
		return nil, err
	}
	_ = file.Close()
	return fileHeader, err
}

// MultipartForm is the parsed multipart form, including file uploads.
func (c *Context) MultipartForm() (*multipart.Form, error) {
	return nil, nil
}

// SaveUploadFile uploads the form file to specific dst.
func (c *Context) SaveUploadFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
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
