package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	lg "github.com/charmbracelet/lipgloss"
	"github.com/sftsrv/tri/theme"
	"github.com/sftsrv/tri/tree"
	"github.com/sftsrv/tri/ui"
)

var usage = lg.JoinVertical(
	lg.Left,

	theme.Primary.MarginBottom(1).Render("tri"),
	"An interactive tree-based search tool with file preview. Give it some path-looking things, and it'll make them readable",

	theme.Secondary.MarginBottom(1).MarginTop(1).Render("Usage"),
	"Pipe in a list of `/` separated stuff, and it'll make them interactive:",
	theme.Faded.MarginBottom(1).MarginTop(1).Render("# files in a directory"),
	"$ find ./ | tri",
	theme.Faded.MarginBottom(1).MarginTop(1).Render("# use alternate preview"),
	"$ find ./ | tri --preview glow",
	theme.Faded.MarginTop(1).Render("# files in a pr"),
	lg.NewStyle().MarginBottom(1).Render("$ git diff --name-only | tri"),
)

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
