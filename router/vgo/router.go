package vgo

import "vgo/context/vgo"

type router struct {
	roots    map[string]*node
	handlers map[string]vgo.HandlerFunc
}

// roots的key类似
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]vgo.HandlerFunc),
	}
}

