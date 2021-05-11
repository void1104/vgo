package core

import (
	"net/http"
	"strings"
	//router2 "vgo/router"
)

// HandlerFunc 定义vgo对于请求的handler
type HandlerFunc func(ctx *Context)

// Engine 实现ServeHTTP方法 -> 实现Handler接口
type Engine struct {
	*GroupRouter // Engine继承了GroupRouter的所有属性和方法，所以*(Engine).engine是指向自己的，将Engine作为最顶层的分组，也就是说Engine拥有Router的所有能力
	router       *Router
	groups       []*GroupRouter
}

// New 引擎的构造方法
func New() (engine *Engine) {
	engine = &Engine{router: newRouter()}
	engine.GroupRouter = &GroupRouter{engine: engine}
	engine.groups = []*GroupRouter{engine.GroupRouter}
	return
}

// addRoute 路由添加方法，调用router模块的方法
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run 定义启动http server的方法
func (engine *Engine) Run(addr string) (err error) {
	// TODO 打印服务器启动日志
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	// 1. 每一次请求都会生成新的context TODO 为请求做缓存
	c := NewContext(w, req)

	c.Handlers = middlewares

	// 2. 交由router的handle函数处理请求
	engine.router.handle(c)
}
