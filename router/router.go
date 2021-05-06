package router

import (
	"log"
	"net/http"
	"vgo/context"
)

type Router struct {
	handlers map[string]HandlerFunc
}

// newRouter 路由构造函数
func newRouter() *Router {
	return &Router{handlers: make(map[string]HandlerFunc)}
}

// AddRoute 静态路由(hash实现)添加函数
func (r *Router) AddRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// handle
func (r *Router) handle(c *context.Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.Status(404)
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
