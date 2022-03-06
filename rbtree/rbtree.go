package rbtree

// An implmentation of the left-leaning red-black 2-3 binary search tree (LLRB BST).
//
// References:
//   https://sedgewick.io/wp-content/themes/sedgewick/papers/2008LLRB.pdf
//   https://algs4.cs.princeton.edu/code/edu/princeton/cs/algs4/RedBlackBST.java

import (
	"golang.org/x/exp/constraints"
)

const (
	COLOR_RED   color = true
	COLOR_BLACK color = false
)

type color bool

type (
	Key   constraints.Ordered
	Value any
)

type Iterator[V any] interface {
	Next() (Iterator[V], bool)
	Value() V
}

type Node[K Key, V Value] struct {
	key    K
	value  V
	left   *Node[K, V]
	right  *Node[K, V]
	parent *Node[K, V]
	color  color
}

func newNode[K Key, V Value](key K, val V, clr color) *Node[K, V] {
	return &Node[K, V]{
		key:   key,
		value: val,
		color: clr,
	}
}

type RedBlackBST[K Key, V Value] struct {
	root *Node[K, V]
}

// New creates an empty instance of a Left-Leaning Red-Black BST.
func New[K Key, V Value]() *RedBlackBST[K, V] {
	return &RedBlackBST[K, V]{}
}

// CompareTo returns > 0 if source is greater than target
func CompareTo[K Key](source, target K) int {
	if source > target {
		return 1
	}
	if source < target {
		return -1
	}

	return 0
}

// Insert a new element
func (t *RedBlackBST[K, V]) Insert(key K, val V) {
	t.root = t.insert(t.root, key, val)
	t.root.color = COLOR_BLACK
}

// insert will recursively traverse down the tree and insert new node at leaf or
// update the value if key exists, then fix by doing rotation or color flip
func (t *RedBlackBST[K, V]) insert(h *Node[K, V], key K, val V) *Node[K, V] {
	if h == nil {
		h = newNode(key, val, COLOR_RED)
		return h
	}

	c := CompareTo(key, h.key)
	switch {
	case c < 0:
		h.left = t.insert(h.left, key, val)
		h.left.parent = h
	case c > 0:
		h.right = t.insert(h.right, key, val)
		h.right.parent = h
	default:
		h.value = val
	}

	// fix height of tree and ensure red links lean left
	if h.right.IsRed() && !h.left.IsRed() {
		h = h.rotateLeft()
	}
	if h.left.IsRed() && h.left.left.IsRed() {
		h = h.rotateRight()
	}
	if h.left.IsRed() && h.right.IsRed() {
		h.flipColors()
	}

	return h
}

// Search by key and returns value, or the zero value of type V if not found
func (t *RedBlackBST[K, V]) Search(key K) V {
	return search(t.root, key)
}

// search does an interative lookup for the key
func search[K Key, V Value](x *Node[K, V], key K) V {
	for x != nil {
		c := CompareTo(key, x.key)
		if c == 0 {
			return x.value
		}

		if c < 0 {
			x = x.left
		} else {
			x = x.right
		}
	}

	// deference and return the zero value based on type
	return *new(V)
}

func (t *RedBlackBST[K, V]) Delete() {
	// TODO
}

func (t *RedBlackBST[K, V]) Begin() Iterator[V] {
	x := t.root
	for x.left != nil {
		x = x.left
	}
	return x
}

func (h *Node[K, V]) Value() V {
	if h == nil {
		return *new(V)
	}

	return h.value
}

func (h *Node[K, V]) Next() (Iterator[V], bool) {
	x := h
	if x.right != nil {
		// one step right
		x = x.right

		// then down to furthest left
		for x.left != nil {
			x = x.left
		}
		return x, true
	}

	// left subtree processed, backtrack up to right only
	p := x.parent
	for x == p.right {
		x = p
		p = p.parent

		if p == nil {
			// all nodes visited, reached up to parent of root which is nil
			return new(Node[K, V]), false
		}
	}

	return p, true
}

// utility functions on Node
func (h *Node[K, V]) IsRed() bool {
	if h == nil {
		return false
	}
	return bool(h.color)
}

func (h *Node[K, V]) rotateLeft() *Node[K, V] {
	x := h.right
	x.parent = h.parent

	h.right = x.left
	if h.right != nil {
		h.right.parent = h
	}

	x.left = h
	if x.left != nil {
		x.left.parent = x
	}

	x.color = x.left.color
	x.left.color = COLOR_RED
	return x
}

func (h *Node[K, V]) rotateRight() *Node[K, V] {
	x := h.left
	h.left = x.right
	if h.left != nil {
		h.left.parent = h
	}
	x.parent = h.parent

	x.right = h
	if x.right != nil {
		x.right.parent = x
	}

	x.color = x.right.color
	x.right.color = COLOR_RED
	return x
}

func (h *Node[K, V]) flipColors() {
	h.color = !h.color
	h.left.color = !h.left.color
	h.right.color = !h.right.color
}
