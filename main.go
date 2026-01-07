package main

import (
	"bufio"
	"os"

	"github.com/sftsrv/tri/tree"
	"github.com/sftsrv/tri/ui"
)

func main() {
	paths := []string{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		paths = append(paths, scanner.Text())
	}

	if len(paths) == 0 {
		panic("Expected to be called with a list of paths from stdin")
	}

	t := tree.PathsToTree(paths)

	ui.Run(t)
}
