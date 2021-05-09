package router

import (
	"fmt"
	"net/http"
	"strings"
	"vgo/context"
	"vgo/utils"
)

/**
前缀树实现动态路由：
	- 参数匹配：例如/p/:lang/doc，可匹配/p/c/doc, /p/go/doc
	- 通配*：例如/static/*filepath, 可匹配/static/fav.ico, /static/js/Query.js
*/

// Router
//	- roots key eg, roots['GET'], roots['POST']
//  - handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
type Router struct {
	roots    map[string]*node       // 存储每种请求方式的Trie树根节点
	handlers map[string]HandlerFunc // 存储每种请求方式的HandlerFunc
	table    utils.Set              // 存储所有注册过的路由
}

// newTrieRouter 前缀树路由构造函数
func newRouter() *Router {
	return &Router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// parsePattern 获取路由字符串的路由数组
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// addRoute 注册路由
func (r *Router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern

	// 判断请求方式是否已注册过路由
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	// 判断路由是否已存在, 存在则panic，将问题暴露给用户
	if r.table.Exist(pattern) {
		panic(fmt.Sprintf("router conflict %s", pattern))
	}

	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// getRoute 路由匹配
func (r *Router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)

	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	// 如果节点不为空，判断是否存在模糊匹配
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

func (r *Router) handle(c *context.Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		// 在调用匹配到的handler前，将解析出来的路由参数赋值给了c.Params，这样就能够在handler中，通过Context对象访问到具体的值了。
		key := c.Method + "-" + n.pattern
		c.Params = params
		c.Handlers = append(c.Handlers, r.handlers[key])
	} else {
		c.Handlers = append(c.Handlers, func(c *context.Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})

		c.Next()
	}
}
