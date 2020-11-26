package core

import (
	"log"
	"net/http"
)

/**
优化点：
	1. 404时返回一个页面，而不是单纯的文字
	2. 当handler传值为null时，使用默认handler
*/

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + " - " + pattern
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {

	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUN: %s\n", c.Path)
	}
}
