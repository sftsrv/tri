package tree

import (
	"fmt"
	"path"
	"slices"
	"strings"
)

type Parts = []string

type kind int

const (
	file kind = iota
	folder
)

type Item struct {
	level int
	name  string
	kind  kind
	tree  *Tree
}

const ICON_FILE = "\uea7b"
const ICON_FOLDER_CLOSED = "\uea83"
const ICON_FOLDER_OPEN = "\uf07c"
const INDENT = "  "

func (s *Item) Expand() {
	s.tree.Expanded = true
}

func (s *Item) ExpandAll() {
	s.tree.ExpandAll()
}

func (t *Tree) ExpandAll() {
	t.Expanded = true
	for _, child := range t.Children {
		child.ExpandAll()
	}
}

func (s *Item) Collapse() {
	s.tree.Expanded = false
}

func (s *Item) CollapseAll() {
	s.tree.CollapseAll()
}

func (t *Tree) CollapseAll() {
	t.Expanded = false
	for _, child := range t.Children {
		child.CollapseAll()
	}
}

func (s Item) IsFile() bool {
	return s.kind == file
}

func (s Item) GetPath() string {
	return path.Join(s.tree.Parts...)
}

func (s Item) icon() string {
	if s.kind == file {
		return ICON_FILE
	}

	if s.tree.Expanded {
		return ICON_FOLDER_OPEN
	}

	return ICON_FOLDER_CLOSED
}

func (t *Tree) Flatten() {
	for childKey, child := range t.Children {
		child.Flatten()
		if len(child.Children) != 1 {
			continue
		}

		grandChildKey := sortedKeys(child.Children)[0]
		grandChild := child.Children[grandChildKey]
		grandChild.Flatten()

		t.Parts = grandChild.Parts
		t.Children[strings.Join(grandChild.Parts, "/")] = grandChild

		delete(child.Children, grandChildKey)
		delete(t.Children, childKey)
	}
}

func (s Item) Render() string {
	return fmt.Sprintf("%s %s %s", strings.Repeat(INDENT, s.level), s.icon(), s.name)
}

func (s Item) Search() string {
	return strings.Join(s.tree.Search(), " ")
}

type Tree struct {
	Expanded bool
	Parts    Parts
	Children map[string]*Tree
}

func (t *Tree) Search() []string {
	tags := t.Parts

	for child, subtree := range t.Children {
		tags = append(tags, child)
		tags = append(tags, subtree.Search()...)
	}

	return tags
}

func newTree(parts Parts) Tree {
	return Tree{
		Expanded: false,
		Parts:    parts,
		Children: map[string]*Tree{},
	}
}

func partsToTreeRec(current Parts, paths []Parts, depth int) *Tree {
	subtrees := map[string][]Parts{}

	for _, parts := range paths {
		if len(parts) > depth {
			segment := parts[depth]
			subtrees[segment] = append(subtrees[segment], parts)
		}
	}

	tree := newTree(current)
	for label, nodes := range subtrees {
		tree.Children[label] = partsToTreeRec(
			append(current, label), nodes, depth+1)
	}

	return &tree
}

func pathsToParts(paths []string) []Parts {
	slices.Sort(paths)

	result := [][]string{}

	for _, path := range paths {
		result = append(result, strings.Split(path, "/"))
	}

	return result
}

func PathsToTree(paths []string) *Tree {
	parts := pathsToParts(paths)
	return partsToTreeRec(Parts{}, parts, 0)
}

func sortedKeys[T any](items map[string]T) []string {
	keys := []string{}

	for key := range items {
		keys = append(keys, key)
	}

	slices.Sort(keys)

	return keys
}

func toItemsRec(tree *Tree, level int) []*Item {
	roots := sortedKeys(tree.Children)

	lines := []*Item{}
	for _, root := range roots {
		children := tree.Children[root]

		kind := file
		if len(children.Children) > 0 {
			kind = folder
		}

		item := &Item{
			level: level,
			name:  root,
			kind:  kind,
			tree:  children,
		}

		lines = append(lines, item)

		if item.tree.Expanded {
			childLines := toItemsRec(children, level+1)
			lines = append(lines, childLines...)
		}

	}

	return lines
}

func ToItems(tree *Tree) []*Item {
	return toItemsRec(tree, 0)
}
