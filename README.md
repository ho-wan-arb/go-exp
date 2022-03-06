# go-exp

Go Experiments to explore generics with data structures.

## Requirements

- Go 1.18rc1 (pending release of 1.18) to use generics.

```sh
go install golang.org/dl/go1.18rc1@latest
go1.18rc1 download
go1.18rc1 test ./...
```

Alias (add to .bashrc to use as default)
```sh
alias go=go1.18rc1
```

## Generics

- Example use of generics from: https://go.dev/doc/tutorial/generics

## LLRB tree

WIP Implementation of a Left-Leaning Red-Black tree with generics.

Left rotation:
```
 [h]           [x]
a   x    ->  h    c
   b c      a b
```

Usage:
```go
tree := New[int, string]()
tree.Insert(1, "a")
v := tree.Search(1)
```
