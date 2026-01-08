# Tri

> Kind of like a nicer `tree --fromfile` command

An interactive tree-based search tool with file preview. Give it some path-looking things, and it'll make them readable

```sh
go install github.com/sftsrv/tri@latest
```

## Usage

Pipe in a list of `/` separated stuff, and it'll make them interactive:

```sh
# files in a directory
find ./ | tri

# files changed in pr
git diff --name-only | tri
```

## TODO ?

- [ ] Multi file select?
- [ ] Custom preview command?
