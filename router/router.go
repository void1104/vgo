package router

import (
	"log"
	"vgo/context"
)

type Router struct {
	handlers map[string]HandlerFunc
}

// newRouter 路由构造函数
func newRouter() *Router {
	return &Router{handlers: make(map[string]HandlerFunc)}
}

// AddRoute 静态路由添加函数
func (r *Router) AddRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// handle
func (r *Router) handle(c *context.Context) {

}
