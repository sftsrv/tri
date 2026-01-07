package tree

import (
	"fmt"
	"slices"
	"strings"
)

type Parts = []string

type Tree struct {
	Children map[string]Tree
}

func newTree() Tree {
	return Tree{
		Children: map[string]Tree{},
	}
}

func partsToTreeRec(paths []Parts, depth int) Tree {
	subtrees := map[string][]Parts{}

	for _, parts := range paths {
		if len(parts) > depth {
			segment := parts[depth]
			subtrees[segment] = append(subtrees[segment], parts)
		}
	}

	tree := newTree()
	for label, nodes := range subtrees {
		tree.Children[label] = partsToTreeRec(nodes, depth+1)
	}

	return tree
}

func pathsToParts(paths []string) []Parts {
	slices.Sort(paths)

	result := [][]string{}

	for _, path := range paths {
		result = append(result, strings.Split(path, "/"))
	}

	return result
}

func PathsToTree(paths []string) Tree {
	parts := pathsToParts(paths)
	return partsToTreeRec(parts, 0)
}

const ICON_FILE = "\uea7b"
const ICON_FOLDER = "\uea83"
const INDENT = "    "

func sortedKeys[T any](items map[string]T) []string {
	keys := []string{}

	for key := range items {
		keys = append(keys, key)
	}

	slices.Sort(keys)

	return keys
}

func printTreeRec(tree Tree, offset int) []string {
	indent := strings.Repeat(INDENT, offset)

	roots := sortedKeys(tree.Children)

	lines := []string{}
	for _, root := range roots {
		children := tree.Children[root]

		icon := ICON_FILE
		if len(children.Children) > 0 {
			icon = ICON_FOLDER
		}

		line := fmt.Sprintf("%s%s %s", indent, icon, root)
		lines = append(lines, line)

		childLines := printTreeRec(children, offset+1)
		lines = append(lines, childLines...)
	}

	return lines
}

func RenderTree(tree Tree) []string {
	return printTreeRec(tree, 0)
}
