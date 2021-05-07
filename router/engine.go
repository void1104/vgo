package router

import (
	"net/http"
	"vgo/context"
)

// HandlerFunc 定义vgo对于请求的handler
type HandlerFunc func(ctx *context.Context)

// Engine 实现了ServeHTTP方法，实现Handler接口
type Engine struct {
	router *Router
}

// New 引擎的构造方法
func New() *Engine {
	return &Engine{router: newRouter()}
}

// addRoute 路由添加方法，调用router模块的方法
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

// GET 定义GET方法的请求
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST 定义POST方法的请求
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run 定义启动http server的方法
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 1. 将req和resp封装为context， 每一次请求都会生成新的context TODO 为请求做缓存
	c := context.NewContext(w, req)
	// 2. 交由router的handle函数处理请求
	engine.router.handle(c)
}
