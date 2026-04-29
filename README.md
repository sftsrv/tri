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
- [ ] Move user search input to separate thread
- [ ] Make splits adjustable via keybinding (resize width of file tree)
- [ ] Make flat mode reactive to searching
- [ ] Make it possible to toggle flat on and off
- [ ] Allow explicit placeholder for file name in output command (like `fzf --preview "cat {}"`)
- [ ] Tests for different formats and structures
- [ ] Fix fuzzy searching
- [ ] Regex based search in editor and via flag?
- [ ] Multi file select?
- [ ] Use `Viewport` (https://github.com/charmbracelet/bubbles) with `reflow` (https://github.com/muesli/reflow) for improved wrapping? - better is being able to provide the size of the active terminal to the underlying process but I have no idea how to do that
