package vgo

import "strings"

// 字典树的节点
type node struct {
	pattern  string  // 节点所匹配的路由，例如/p/:lang
	part     string  // 路由中的一部分，例如:lang
	children []*node // 子节点，也就是下一级路由
	isWild   bool    // 是否模糊匹配，part含有 : 或 * 时为true
}

// 第一个匹配到的节点，用于插入（事实上，一颗tire树中也不会有内容重复的节点）
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配到的节点，用于搜索
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 注册路由（将配置的url插入tire树）
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height { // 如果树的高度等于数组长度，说明已到最后一级路由
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil { // 在遍历tire树插入的过程中，如果遇到空节点，则一路新建
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") { // 如果到达最后一级节点，或者当级路由以*为前缀
		if n.pattern == "" { // 如果当前节点只是别的路由的子集，则返回
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	// dfs遍历tire树，找到对应路由节点
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil { // 如果不为空，则搜索到了路径
			return result
		}
	}

	return nil
}
