# Tri

> Kind of like an interactive `tree --fromfile` command

An interactive tree-based search tool with file preview. Give it some path-looking things, and it'll make them readable

```sh
go install github.com/sftsrv/tri@latest
```

![Screen recording of Tri in action](./images/tri.gif)

## Usage

```
# files in a directory
find ./ | tri

# use alternate preview
find ./ | tri --preview glow

# files in a pr
git diff --name-only | tri
```

### Using Patterns

Use Regexps to parse the input string to create more complex commands

> Escaping of regex special chars will depend on your shell

```
# viewing a formatted git log and showing the changes for each hash
git log --pretty=format:"%h %f"
| tri --preview "git show $1" --pattern "^(\w+)"

# using a more complex regexp and command
git log --pretty=format:"%h %f"
| tri --preview "echo title: $title hash: $hash" --pattern "(?<hash>\w+) (?<title>.*)"
```

### Using Regexps

The structure of the regexp provided should be compatible with Go's implementation,
with the following affordances made:

1. Named capture groups can be specified as `<name>` instead of `?P<name>`
2. Regexps are always matched against a single line - so multiline captures are not meaningful

### Using Patterns

References in pattern are indicated with a `$` in the pattern

        - `$` refers to the entire match
        - `$0` (entire match), `$1` (first capture group), etc. refer to capture groups in order of capture
        - `$<name>` or `$name` refer to named capture groups

## Features / Ideas / TODOs

There are a lot of smaller improvements still left, but I've been using it for some time now and seems to work fine for me - but if you're keen to pick something up then do feel free to

- [x] Search
- [x] Preview with syntax highlighting (using bat if available)
- [x] Custom preview command
- [x] File selection
- [x] Expand/Collapse folders
- [x] Flatten direct paths (using `--flat` flag)
- [x] Print tree (using `--print` flag)
- [x] Async/non-blocking previews
- [x] Make splits adjustable via keybinding (resize width of file tree)
- [x] Allow explicit placeholder for file name in output command (like how it works for)
  - Uses Regexp for pattern definition and dynamic commands
- [ ] Move user search input to separate thread
- [ ] Make flat mode reactive to searching
- [ ] Make it possible to toggle flat on and off
- [ ] Tests for different formats and structures
- [ ] Fix fuzzy searching
- [ ] Regex based search in editor and via flag?
- [ ] Multi file select?
- [ ] Use `Viewport` (https://github.com/charmbracelet/bubbles) with `reflow` (https://github.com/muesli/reflow) for improved wrapping? - better is being able to provide the size of the active terminal to the underlying process but I have no idea how to do that
