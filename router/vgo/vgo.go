package vgo

import "net/http"

/**
框架入口
*/

// 定义vgo使用的请求处理程序
type HandlerFunc func(*Context)

// 实现ServeHTTP方法，作为net/http包处理请求的入口
type Engine struct {
	router *router
}

// Engine的构造方法
func New() *Engine {
	return &Engine{router: newRouter()}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

// HTTP - GET方法
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// HTTP - POST方法
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// 启动http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

/**
实现Serve接口，接管所有的HTTP请求
- HTTP请求是并发的，但每个请求都会调用ServeHTTP，这个方法中，
  每次都会创建新的context，不会对同一个context进行写入。
*/
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
