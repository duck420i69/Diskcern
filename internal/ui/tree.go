package ui

import (
	"sort"
	"strings"
	"time"

	"gioui.org/widget"
	"github.com/diskcern/diskcern/internal/models"
)

type TreeNode struct {
	Name      string
	Path      string
	Size      int64
	IsDir     bool
	Expanded  bool
	Children  map[string]*TreeNode
	
	SortedChildren []*TreeNode
	
	Clickable widget.Clickable
	LastClick time.Time
}

func BuildTree(records []models.FileRecord) *TreeNode {
	root := &TreeNode{
		Name:     "Root",
		Path:     "",
		IsDir:    true,
		Expanded: true,
		Children: make(map[string]*TreeNode),
	}

	for _, r := range records {
		parts := strings.Split(r.Path, "\\")
		
		current := root
		currentPath := ""
		
		for i, part := range parts {
			if part == "" {
				continue
			}
			
			if currentPath == "" {
				currentPath = part
			} else {
				currentPath = currentPath + "\\" + part
			}
			
			if _, exists := current.Children[part]; !exists {
				isFolder := true
				if i == len(parts)-1 {
					isFolder = r.IsDir
				}
				
				current.Children[part] = &TreeNode{
					Name:     part,
					Path:     currentPath,
					IsDir:    isFolder,
					Children: make(map[string]*TreeNode),
				}
			}
			
			current = current.Children[part]
			
			if i == len(parts)-1 {
				current.Size = r.Size
			}
		}
	}
	
	// Since WizTree already calculates sizes for folders, we might not need to accumulate.
	// However, if we filter out provider-matched folders, the parent folder sizes won't be accurate anymore!
	// Re-accumulate sizes recursively to fix this!
	recalculateSizes(root)
	
	sortTree(root)
	
	return root
}

func recalculateSizes(node *TreeNode) int64 {
	var total int64
	for _, child := range node.Children {
		total += recalculateSizes(child)
	}
	// If it's a file, keep its original size. If folder, its size is sum of children.
	// But what if it's a folder that contains files NOT captured? 
	// WizTree records every file, so recalculating from leaf nodes is perfectly accurate.
	if !node.IsDir || len(node.Children) == 0 {
		return node.Size
	}
	node.Size = total
	return total
}

func sortTree(node *TreeNode) {
	node.SortedChildren = make([]*TreeNode, 0, len(node.Children))
	for _, child := range node.Children {
		node.SortedChildren = append(node.SortedChildren, child)
		sortTree(child)
	}
	
	sort.Slice(node.SortedChildren, func(i, j int) bool {
		return node.SortedChildren[i].Size > node.SortedChildren[j].Size
	})
}

func FlattenTree(root *TreeNode, scores map[string]int) []*TreeNode {
	var result []*TreeNode
	
	var traverse func(node *TreeNode)
	traverse = func(node *TreeNode) {
		if node != root && scores[node.Path] >= 0 {
			result = append(result, node)
		}
		if node.Expanded || node == root {
			for _, child := range node.SortedChildren {
				if scores[child.Path] >= 0 {
					traverse(child)
				}
			}
		}
	}
	traverse(root)
	
	return result
}
