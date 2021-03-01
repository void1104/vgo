package vgo

import (
	"fmt"
	"net/http"
)

/**
提供给框架用户，用来定义路由映射的处理方法
*/
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type Engine struct {
	router map[string]HandlerFunc // 路由映射表
}

// 实例化函数
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// 所有请求方法的基石函数
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// 启动http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

/**
解析请求的路径，查找路由映射表，如果查到，就执行注册的处理方法。如果查不到，就返回404 NOT FOUND
*/
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		w.WriteHeader(http.StatusNotFound) // 修改状态码
		fmt.Fprintf(w, "404 NOT FOUND:%s\n", req.URL)
	}
}
