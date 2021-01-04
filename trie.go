/*
 * @Description:
 * @Version: 1.0
 * @Author: Curricane
 * @Date: 2020-12-30 15:33:55
 * @LastEditors: Curricane
 * @LastEditTime: 2020-12-30 16:24:34
 * @FilePath: /golang/cweb/trie.go
 * @Copyright (C) 2020 Curricane. All rights reserved.
 */

package cweb

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s}, part=%s, isWild=%t", n.pattern, n.part, n.isWild)
}

// insert 插入节点 {pattern}完整url路径，{parts}url路径的每一段，{height}当前插入的高度
func (n *node) insert(pattern string, parts []string, height int) {
	// 已经插入到最深，不需要继续插入，节点写入pattern成员
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height] // 当前要匹配的一段路径
	child := n.matchChild((part))
	if child == nil { // n没有匹配的child，n节点为插入点
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}

	// 继续插入
	child.insert(pattern, parts, height+1)
}

// 查询匹配路径的节点
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

// matchChild 第一个匹配part成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// travel 深度遍历所有节点，找到所有的路由路径
func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}

	for _, child := range n.children {
		child.travel(list)
	}
}

// matchChildren 所有匹配part成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
