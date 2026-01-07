package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sftsrv/tri/tree"
)

type Parts = []string

func main() {
	paths := []string{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		paths = append(paths, scanner.Text())
	}

	if len(paths) == 0 {
		panic("Expected to be called with a list of paths from stdin")
	}

	items := tree.PathsToTree(paths)
	lines := tree.RenderTree(items)

	for _, line := range lines {
		fmt.Println(line)
	}
}
