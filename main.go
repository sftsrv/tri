package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sftsrv/tri/tree"
	"github.com/sftsrv/tri/ui"
)

var usage = `tri

An interactive tree-based search tool with file preview. Give it some path-looking things, and it'll make them readable

## Usage

'''
# files in a directory
find ./ | tri

# use alternate preview
find ./ | tri --preview glow

# files in a pr
git diff --name-only | tri
'''

### Using Patterns

Use Regexps to parse the input string to create more complex commands

> Escaping of regex special chars will depend on your shell

'''
# viewing a formatted git log and showing the changes for each hash
git log --pretty=format:"%h %f"
| tri --preview "git show $1" --pattern "^(\w+)"

# using a more complex regexp and command
git log --pretty=format:"%h %f"
| tri --preview "echo title: $title hash: $hash" --pattern "(?<hash>\w+) (?<title>.*)"
'''

### Using Regexps

The structure of the regexp provided should be compatible with Go's implementation,
with the following affordances made:

1. Named capture groups can be specified as '<name>' instead of '?P<name>'
2. Regexps are always matched against a single line - so multiline captures are not meaningful

### Using Patterns

References in pattern are indicated with a '$' in the pattern

        - '$' refers to the entire match
        - '$0' (entire match), '$1' (first capture group), etc. refer to capture groups in order of capture
        - '$<name>' or '$name' refer to named capture groups

`

func main() {
	help := flag.Bool("help", false, "show help menu")
	print := flag.Bool("print", false, "print tree (non interactive)")
	flat := flag.Bool("flat", false, "flatten direct paths")
	preview := flag.String("preview", "", "command to use for file preview")
	pattern := flag.String("pattern", "", "pattern to use when parsing path")

	flag.Parse()

	if *help {
		fmt.Println(usage)
		flag.Usage()
		return
	}

	paths := []string{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		paths = append(paths, line)
	}

	if len(paths) == 0 {
		panic("Expected to be called with a list of paths from stdin")
	}

	t := tree.PathsToTree(paths)

	if *print {
		t.ExpandAll()
		if *flat {
			t.Flatten()
		}
		fmt.Println(tree.Render(t))
		return
	}

	ui.Run(t, *preview, *pattern, *flat)
}
