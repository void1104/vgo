package router

import "strings"

// node 前缀树的节点
type node struct {
	pattern  string  // 待匹配路由，例如/p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如[doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，part含有 : 或 * 时为true
}

// matchChild 返回第一个匹配成功的节点（精确匹配 or 模糊匹配），用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// matchChildren 返回所有匹配成功的节点（精确匹配 or 模糊匹配），用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// insert 插入节点 TODO 路由冲突时，将问题暴露给用户
func (n *node) insert(pattern string, parts []string, height int) {
	// 如果走到叶子节点，并且路由也走到最后的分隔符，则把当前节点路径注册为路由
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)

	// 如果没找到节点，则顺着路径生成一个空节点
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}

	// 如果找到了，则递归进入子节点
	child.insert(pattern, parts, height+1)
}

// search 查找节点
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	// dfs遍历所有子节点，直到找到对应节点为止
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
