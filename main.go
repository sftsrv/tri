package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/sftsrv/tri/tree"
	"github.com/sftsrv/tri/ui"
)

var usage = `tri

An interactive tree-based search tool with file preview. Give it some path-looking things, and it'll make them readable

Usage

# files in a directory
$ find ./ | tri
# use alternate preview
$ find ./ | tri --preview glow
# files in a pr
$ git diff --name-only | tri

`

func main() {
	help := flag.Bool("help", false, "show help menu")
	preview := flag.String("preview", "", "command to use for file preview")

	flag.Parse()

	if *help {
		fmt.Println(usage)
		flag.Usage()
		return
	}

	paths := []string{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		paths = append(paths, scanner.Text())
	}

	if len(paths) == 0 {
		panic("Expected to be called with a list of paths from stdin")
	}

	t := tree.PathsToTree(paths)

	ui.Run(t, *preview)
}
