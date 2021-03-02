package vgo

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
	isWild   bool    // 是否精确匹配
}

// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
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
	if len(parts) == height {
		n.pattern = pattern
		return
	}

}