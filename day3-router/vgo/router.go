package vgo

import (
	"net/http"
	"strings"
)

type Router struct {
	roots    map[string]*node // 每种http方法，对应一颗tire树
	handlers map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

/**
  tips:只有*的路由也是允许的
  将路由url根据"/"分割，转换为字符串数组
*/
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' { // 如果遇到*号，就不用继续分割了
				break
			}
		}
	}
	return parts
}

// 注册路由
func (r *Router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern // 以方法和url作为key，对应的handler作为value
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0) // 找到http方法对应的tire树，并插入
	r.handlers[key] = handler
}

// 寻找路由
func (r *Router) getRouter(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

// 处理函数
func (r *Router) handle(c *Context) {
	n, params := r.getRouter(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND:%s\n", c.Path)
	}
}
