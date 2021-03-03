package vgo

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

/**
roots的key类似 roots['GET'] roots['POST']
handlers的key类似 handlers['GET -/p/:lang/doc'], handlers['POST-/p/book']
*/

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

/**
将路由切分，只有一个*的路由也是允许的
*/
func parsePattern(pattern string) (parts []string) {
	vs := strings.Split(pattern, "/")
	parts = make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			// 如果以*开头，则不需在意后续路由，*表示全量匹配
			if item[0] == '*' {
				break
			}
		}
	}
	return
}

/**
添加路由
*/
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	// 每个HTTP方法，有各自的Trie树
	_, ok := r.roots[method]
	// 如果当前HTTP方法不存在路由，则创建根节点
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

/**
路由匹配
*/
func (r *router) getRoute(method string, path string) (n *node, params map[string]string) {
	searchParts := parsePattern(path)
	params = make(map[string]string)
	// 根据HTTP方法找到对应的路由
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	n = root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			// 如果为模糊匹配的路由，则将截取的路由存入上下文中，方便后续获得
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

/**
请求入口处理函数

*/
func (r *router) handle(c *Context) {
	// 确定路由匹配后，调用响应的handlerFunc执行
	n, params := r.getRoute(c.Req.Method, c.Req.URL.Path)
	if n != nil {
		c.Params = params
		key := c.Req.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Req.URL.Path)
	}
}
