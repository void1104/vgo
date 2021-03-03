package vgo

import "strings"

/**
Trie树实现动态路由
	- 参数匹配：例如 /p/:lang/doc, 可以匹配/p/c/doc/和/p/go/doc
	- 通配*: 例如 /static/*filepath, 可以匹配/static/fav.ico，可以匹配/static/js/jQuery.js，
			这种模式常用于静态服务器，能够递归匹配子路径
*/

type node struct {
	pattern  string  // 待匹配路由，例如/p/:lang
	part     string  // 路由中的一部分，例如:lang
	children []*node // 子节点，例如[doc, tutorial, intro]
	isWild   bool    // 是否模糊匹配
}

// 第一个匹配成功的节点，用于插入
func (n *node) matchFirstChild(part string) *node {
	// 遍历所有子节点，如果匹配成功或者子节点为模糊匹配，返回子节点
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

/**
注册路由规则，映射handler
*/
func (n *node) insert(pattern string, parts []string, height int) {
	// 到传入路由的最后一级，将当前节点的路由匹配设置为传入pattern
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	//  匹配当前part的子节点
	child := n.matchFirstChild(part)
	// 如果没有匹配到当前part的节点，则新建一个
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

/**
匹配路由规则，查找到对应的handler
*/
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
