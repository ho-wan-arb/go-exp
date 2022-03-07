package rbtree

// An implmentation of the left-leaning red-black 2-3 binary search tree (LLRB BST).
//
// References:
//   https://sedgewick.io/wp-content/themes/sedgewick/papers/2008LLRB.pdf
//   https://algs4.cs.princeton.edu/code/edu/princeton/cs/algs4/RedBlackBST.java

import (
	"fmt"
	"strings"

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

type Iterator[K Key, V Value] struct {
	tree    *RBTree[K, V]
	current *RBNode[K, V]
}

type RBNode[K Key, V Value] struct {
	key    K
	value  V
	left   *RBNode[K, V]
	right  *RBNode[K, V]
	parent *RBNode[K, V]
	color  color
}

func newNode[K Key, V Value](key K, val V, clr color) *RBNode[K, V] {
	return &RBNode[K, V]{
		key:   key,
		value: val,
		color: clr,
	}
}

type RBTree[K Key, V Value] struct {
	root *RBNode[K, V]
}

// New creates an empty instance of a Left-Leaning Red-Black BST.
func New[K Key, V Value]() *RBTree[K, V] {
	return &RBTree[K, V]{}
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
func (t *RBTree[K, V]) Insert(key K, val V) {
	t.root = t.insert(t.root, key, val)
	t.root.color = COLOR_BLACK
}

// insert will recursively traverse down the tree and insert new node at leaf or
// update the value if key exists, then fix by doing rotation or color flip
func (t *RBTree[K, V]) insert(cur *RBNode[K, V], key K, val V) *RBNode[K, V] {
	if cur == nil {
		cur = newNode(key, val, COLOR_RED)
		return cur
	}

	c := CompareTo(key, cur.key)
	switch {
	case c < 0:
		cur.left = t.insert(cur.left, key, val)
		cur.left.parent = cur
	case c > 0:
		cur.right = t.insert(cur.right, key, val)
		cur.right.parent = cur
	default:
		cur.value = val
	}

	// fix height of tree and ensure red links lean left
	if cur.right.isRed() && !cur.left.isRed() {
		cur = cur.rotateLeft()
	}
	if cur.left.isRed() && cur.left.left.isRed() {
		cur = cur.rotateRight()
	}
	if cur.left.isRed() && cur.right.isRed() {
		cur.flipColors()
	}

	return cur
}

// Search by key and returns value, or the zero value of type V if not found
func (t *RBTree[K, V]) Search(key K) (V, bool) {
	cur := t.root
	for cur != nil {
		c := CompareTo(key, cur.key)
		if c == 0 {
			return cur.value, true
		}

		if c < 0 {
			cur = cur.left
		} else {
			cur = cur.right
		}
	}

	// deference and return zero value
	return *new(V), false
}

func (t *RBTree[K, V]) Delete() {
	// TODO
}

func (t *RBTree[K, V]) Begin() *Iterator[K, V] {
	cur := t.root
	for cur.left != nil {
		cur = cur.left
	}

	return &Iterator[K, V]{
		tree:    t,
		current: cur,
	}
}

// End moves to one past the last element.
func (t *RBTree[K, V]) Last() *Iterator[K, V] {
	cur := t.root
	for cur.right != nil {
		cur = cur.right
	}

	return &Iterator[K, V]{
		tree:    t,
		current: cur,
	}
}

// End moves to one past the last element.
func (t *RBTree[K, V]) End() *Iterator[K, V] {
	return &Iterator[K, V]{
		tree:    t,
		current: nil,
	}
}

// String prints the tree in a visual format row by row.
func (rb *RBTree[K, V]) String() string {
	d := 0
	list := map[int][]string{}
	traverseByDepth(rb.root, d, list)

	sb := strings.Builder{}
	sb.WriteString("----\n")
	for i := 1; i <= len(list); i++ {
		sb.WriteString(fmt.Sprintf("[depth %d]:  ", i))
		sb.WriteString(fmt.Sprintf("%v\n", strings.Join(list[i], " | ")))
	}

	return sb.String()
}

// Next does an in-order traversal through a binary search tree.
func (it *Iterator[K, V]) Next() bool {
	cur := it.current
	if cur == nil {
		begin := it.tree.Begin()
		it.current = begin.current
		return true
	}

	if cur.right != nil {
		// one step right
		cur = cur.right

		// then down to furthest left
		for cur.left != nil {
			cur = cur.left
		}
		it.current = cur
		return true
	}

	// left subtree processed, backtrack up to right only
	for cur == cur.parent.right {
		cur = cur.parent

		if cur.parent == nil {
			// all nodes visited, reached up to parent of root which is nil
			it.current = nil
			return false
		}
	}

	it.current = cur.parent
	return true
}

// Next does an in-order traversal through a binary search tree in reverse.
func (it *Iterator[K, V]) Prev() bool {
	cur := it.current
	if cur == nil {
		begin := it.tree.Last()
		it.current = begin.current
		return true
	}

	if cur.left != nil {
		// one step right
		cur = cur.left

		// then down to furthest right
		for cur.right != nil {
			cur = cur.right
		}
		it.current = cur
		return true
	}

	// right subtree processed, backtrack up to left only
	for cur == cur.parent.left {
		cur = cur.parent

		if cur.parent == nil {
			// all nodes visited, reached up to parent of root which is nil
			it.current = nil
			return false
		}
	}

	it.current = cur.parent
	return true
}

func (it *Iterator[K, V]) Key() K {
	if it.current == nil {
		// deference and return zero value
		return *new(K)
	}

	return it.current.key
}

func (it *Iterator[K, V]) Value() V {
	if it.current == nil {
		// deference and return zero value
		return *new(V)
	}

	return it.current.value
}

func (n *RBNode[K, V]) isRed() bool {
	if n == nil {
		return false
	}
	return bool(n.color)
}

func (n *RBNode[K, V]) rotateLeft() *RBNode[K, V] {
	cur := n.right
	cur.parent = n.parent

	n.right = cur.left
	if n.right != nil {
		n.right.parent = n
	}

	cur.left = n
	if cur.left != nil {
		cur.left.parent = cur
	}

	cur.color = cur.left.color
	cur.left.color = COLOR_RED
	return cur
}

func (n *RBNode[K, V]) rotateRight() *RBNode[K, V] {
	cur := n.left
	n.left = cur.right
	if n.left != nil {
		n.left.parent = n
	}
	cur.parent = n.parent

	cur.right = n
	if cur.right != nil {
		cur.right.parent = cur
	}

	cur.color = cur.right.color
	cur.right.color = COLOR_RED
	return cur
}

func (n *RBNode[K, V]) flipColors() {
	n.color = !n.color
	n.left.color = !n.left.color
	n.right.color = !n.right.color
}

func traverseByDepth[K Key, V Value](cur *RBNode[K, V], d int, list map[int][]string) {
	if cur == nil {
		return
	}

	curKey := fmt.Sprintf("%v", cur.key)

	if !cur.isRed() {
		d++
		list[d] = append(list[d], curKey)
	} else {
		// join 2 nodes: red link should lean left, so smaller number should always be in front
		list[d][len(list[d])-1] = fmt.Sprintf("(%v,%v)", curKey, list[d][len(list[d])-1])
	}

	traverseByDepth(cur.left, d, list)
	traverseByDepth(cur.right, d, list)
}
