package main

import (
	"bufio"
	"fmt"
	"os"
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

func partsToTree(paths []Parts) Tree {
	return partsToTreeRec(paths, 0)
}

const ICON_FILE = "\uea7b"
const ICON_FOLDER = "\uea83"

func sortedKeys[T any](items map[string]T) []string {
	keys := []string{}

	for key := range items {
		keys = append(keys, key)
	}

	slices.Sort(keys)

	return keys
}

func printTreeRec(tree Tree, offset int) {
	indent := strings.Repeat("  ", offset)

	roots := sortedKeys(tree.Children)

	for _, root := range roots {
		children := tree.Children[root]

		icon := ICON_FILE
		if len(children.Children) > 0 {
			icon = ICON_FOLDER
		}

		fmt.Printf("%s%s %s\n", indent, icon, root)
		printTreeRec(children, offset+1)
	}
}

func printTree(tree Tree) {
	printTreeRec(tree, 0)
}

func pathsToParts(paths []string) []Parts {
	slices.Sort(paths)

	result := [][]string{}

	for _, path := range paths {
		result = append(result, strings.Split(path, "/"))
	}

	return result
}

func main() {
	paths := []string{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		paths = append(paths, scanner.Text())
	}

	if len(paths) == 0 {
		panic("Expected to be called with a list of paths from stdin")
	}

	parts := pathsToParts(paths)
	tree := partsToTree(parts)

	printTree(tree)
}
