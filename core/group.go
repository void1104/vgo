package core

import (
	"log"
)

// GroupRouter Group对象需要有访问Router的能力，为了方便，
// 我们可以在Group中，保存一个指针，指向Engine，整个框架的所
// 有资源都是由Engine统一协调的，那么就可以通过Engine间接的
// 访问各种接口了
type GroupRouter struct {
	prefix      string        // 分组的前缀
	middlewares []HandlerFunc // 分组对应的中间件
	parent      *GroupRouter  // 支持嵌套
	engine      *Engine       // 所有分组共享一个Engine实例，
}

// Group 创建一个新的路由分组，记住所有分组共享一个Engine
func (group *GroupRouter) Group(prefix string) (gr *GroupRouter) {
	engine := group.engine
	gr = &GroupRouter{
		prefix: group.prefix + prefix,
		parent: group,
		engine: group.engine,
	}
	engine.groups = append(engine.groups, gr)
	return
}

// Use 为路由组注册中间件
func (group *GroupRouter) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// addRoute 分组添加路由
func (group *GroupRouter) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *GroupRouter) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *GroupRouter) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}
