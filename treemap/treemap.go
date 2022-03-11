package treemap

// An implmentation of a treemap backed by a left-leaning red-black 2-3 binary search tree (LLRB BST).
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
	red   color = true
	black color = false
)

type (
	key   constraints.Ordered
	val   any
	color bool
)

// TreeMap holds a binary search tree that can be used as a map.
type TreeMap[K key, V val] struct {
	root *node[K, V]
}

// New creates an empty instance of a TreeMap.
func New[K key, V val]() *TreeMap[K, V] {
	return &TreeMap[K, V]{}
}

type node[K key, V val] struct {
	key    K
	value  V
	color  color
	left   *node[K, V]
	right  *node[K, V]
	parent *node[K, V]
}

func newNode[K key, V val](k K, v V, c color) *node[K, V] {
	return &node[K, V]{
		key:   k,
		value: v,
		color: c,
	}
}

// CompareTo returns > 0 if source is greater than target
func CompareTo[K key](source, target K) int {
	if source > target {
		return 1
	}
	if source < target {
		return -1
	}

	return 0
}

// Insert a new element
func (t *TreeMap[K, V]) Insert(key K, val V) {
	t.root = t.insert(t.root, key, val)
	t.root.color = black
}

// insert will recursively traverse down the tree and insert new node at leaf or
// update the value if key exists, then fix by doing rotation or color flip
func (t *TreeMap[K, V]) insert(cur *node[K, V], key K, val V) *node[K, V] {
	if cur == nil {
		cur = newNode(key, val, red)
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
func (t *TreeMap[K, V]) Search(key K) (V, bool) {
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

// Begin moves iterator in front of first element.
func (t *TreeMap[K, V]) Begin() *Iterator[K, V] {
	cur := t.root
	for cur.left != nil {
		cur = cur.left
	}

	return &Iterator[K, V]{
		tree:    t,
		current: cur,
	}
}

// Last moves iterator in front of the last element.
func (t *TreeMap[K, V]) Last() *Iterator[K, V] {
	cur := t.root
	for cur.right != nil {
		cur = cur.right
	}

	return &Iterator[K, V]{
		tree:    t,
		current: cur,
	}
}

// End moves iterator to behind the last element.
func (t *TreeMap[K, V]) End() *Iterator[K, V] {
	return &Iterator[K, V]{
		tree:    t,
		current: nil,
	}
}

// String prints the tree in a visual format row by row.
func (t *TreeMap[K, V]) String() string {
	d := 0
	list := map[int][]string{}
	traverseByDepth(t.root, d, list)

	sb := strings.Builder{}
	sb.WriteString("----\n")
	for i := 1; i <= len(list); i++ {
		sb.WriteString(fmt.Sprintf("[depth %d]:  ", i))
		sb.WriteString(fmt.Sprintf("%v\n", strings.Join(list[i], " | ")))
	}

	return sb.String()
}

type Iterator[K key, V val] struct {
	tree    *TreeMap[K, V]
	current *node[K, V]
}

// Next does an in-order traversal through a binary search tree.
func (it *Iterator[K, V]) Next() bool {
	cur := it.current
	if cur == nil {
		begin := it.tree.Begin()
		it.current = begin.current
		return true
	}

	// try to go one step right then all the way to left
	if cur.right != nil {
		cur = cur.right

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

// Prev does an in-order traversal through a binary search tree in reverse.
func (it *Iterator[K, V]) Prev() bool {
	// steps are just the reverse of forward traversal
	cur := it.current
	if cur == nil {
		begin := it.tree.Last()
		it.current = begin.current
		return true
	}

	if cur.left != nil {
		cur = cur.left

		for cur.right != nil {
			cur = cur.right
		}
		it.current = cur
		return true
	}

	for cur == cur.parent.left {
		cur = cur.parent

		if cur.parent == nil {
			it.current = nil
			return false
		}
	}

	it.current = cur.parent
	return true
}

// Key returns the key at the current position of iterator and returns the zero value if nil.
func (it *Iterator[K, V]) Key() K {
	if it.current == nil {
		return *new(K)
	}

	return it.current.key
}

// Key returns the key at the current position of iterator and returns the zero value if nil.
func (it *Iterator[K, V]) Value() V {
	if it.current == nil {
		return *new(V)
	}

	return it.current.value
}

func (n *node[K, V]) isRed() bool {
	if n == nil {
		return false
	}
	return bool(n.color)
}

func (n *node[K, V]) rotateLeft() *node[K, V] {
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
	cur.left.color = red
	return cur
}

func (n *node[K, V]) rotateRight() *node[K, V] {
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
	cur.right.color = red
	return cur
}

func (n *node[K, V]) flipColors() {
	n.color = !n.color
	n.left.color = !n.left.color
	n.right.color = !n.right.color
}

func traverseByDepth[K key, V val](cur *node[K, V], d int, list map[int][]string) {
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