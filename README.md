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

Right rotation:
```
   [h]         [x]
 x    c  ->   a   h
a b              b c
```

API
```go
type (
	Key   constraints.Ordered
	Value any
)

type Iterator[V any] interface {
	Next() (Iterator[V], bool)
	Value() V
}

type RedBlackBST[K Key, V Value]
   func New() *RedBlackBST
   func (t *RedBlackBST) Insert(key K, val V)
   func (t *RedBlackBST) Search(key K) V
   func (t *RedBlackBST) Begin() Iterator[V]

type Node[K Key, V Value]
   func (h *Node[K, V]) Value() V
   func (h *Node[K, V]) Next() (Iterator[V], bool)
   func (h *Node[K, V]) IsRed() bool
```
