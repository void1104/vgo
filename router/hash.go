package router

import (
	"log"
	"net/http"
	"vgo/context"
)

/**
hash表实现静态路由
*/

type hashRouter struct {
	handlers map[string]HandlerFunc
}

// newRouter 路由构造函数
func newHashRouter() *hashRouter {
	return &hashRouter{handlers: make(map[string]HandlerFunc)}
}

// AddRoute 静态路由(hash实现)添加函数
func (r *hashRouter) AddRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// handle
func (r *hashRouter) handle(c *context.Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.Status(404)
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
