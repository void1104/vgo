package vgo

import (
	"log"
	"net/http"
)

/**
将路由相关的方法和结构提取出来，放在这里。
方便我们下一次对router的功能进行增强，例如提供动态路由的支持。
*/
type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

/**
将路由以及对应的HandlerFunc加入handlers中
*/
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	key := c.Req.Method + "-" + c.Req.URL.Path
	if handler, ok := r.handlers[key]; ok {
		// 通过key找到注册的HandlerFunc,用它去执行具体流程
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND:%s\n", c.Req.URL.Path)
	}
}
