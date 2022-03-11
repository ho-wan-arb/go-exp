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

## TreeMap

A Generic treemap backed by a balanced Binary Search Tree (BST).
The treemap can be traversed in sorted order of keys using the Iterator.
A custom comparator function can be used when initializing the treemap.

This uses a Left-Leaning 2-3 Red-Black tree.

Key properties:
- Red links can 'glue' 2 nodes together, which have 3 children (hence 2-3 tree).
- Red links have to lean left, cannot have 2 or more red links together for a 2-3 tree.
- Depth to all nodes is the same when counting blank links only.

Complexity (in general for balanced BST):
- O(logN) to search
- O(logN) to insert
- O(logN) to delete

### Details

Left rotation:
```
 [b]           [c]
a   c    ->  b    e
   d e      a d
```

Right rotation:
```
   [c]         [b]
 b    e  ->   a   c
a d              d e
```
