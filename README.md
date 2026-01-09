# Tri

> Kind of like an interactive `tree --fromfile` command

An interactive tree-based search tool with file preview. Give it some path-looking things, and it'll make them readable

```sh
go install github.com/sftsrv/tri@latest
```

![Screen recording of Tri in action](./images/tri.gif)

## Usage

Pipe in a list of `/` separated stuff, and it'll make them interactive:

```sh
# files in a directory
find ./ | tri

# files changed in pr
git diff --name-only | tri

# using a custom previewer
git diff --name-only | tri --preview "git diff HEAD --"
```

## Features / TODOs

- [x] Search
- [x] Preview with syntax highlighting (using bat if available)
- [x] Custom preview command
- [x] File selection
- [x] Expand/Collapse folders
- [x] Flatten direct paths (using `--flat` flag)
- [x] Print tree (using `--print` flag)
- [ ] Regex based search in editor and via flag
- [ ] Tests for different formats and structures
- [ ] Fix fuzzy searching
- [ ] Flatten files during search in flat mode
- [ ] Multi file select?
