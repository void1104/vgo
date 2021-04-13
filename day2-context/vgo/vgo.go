package vgo

import "net/http"

/**
重点改进:
	- 将路由(router)独立初来, 方便之后增强.
	- 设计上下文(context), 封装request和response,提供对JSON,HTML等返回类型的支持.
*/

// HandlerFunc defines the request handler used by vgo
type HandlerFunc func(*Context)

// Engine implements the interface of ServeHTTP
type Engine struct {
	router *router
}

// New is the constructor of vgo.Engine
func New() *Engine {
	return &Engine{router: newRouter()}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

/**
作为所有请求的入口, 每次请求都会生成一个新的Context
 */
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
