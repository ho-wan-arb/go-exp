// Package treemap implements a generic treemap backed by a Binary Search Tree.
// The treemap can be traversed in sorted order of keys using the Iterator.
// A custom comparator function can be used when initializing the treemap.
// Generics require go version > 1.18 to be used.
//
// A Left-Leaning 2-3 Red-Black (LLRB) tree is used as the self-balancing Binary Search Tree (BST).
// References:
//   https://sedgewick.io/wp-content/themes/sedgewick/papers/2008LLRB.pdf
//   https://algs4.cs.princeton.edu/code/edu/princeton/cs/algs4/RedBlackBST.java
package treemap

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/exp/constraints"
)

const (
	red   color = true
	black color = false
)

type (
	key   any
	val   any
	color bool
)

// TreeMap is a balanced binary search tree that can be used as a map.
type TreeMap[K key, V val] struct {
	root       *node[K, V]
	comparator Comparator[K]
	length     int
}

// NewWithComparator creates a new TreeMap using the default comparator (< and >).
func New[K constraints.Ordered, V val]() *TreeMap[K, V] {
	return &TreeMap[K, V]{
		comparator: defaultComparator[K],
	}
}

// NewWithComparator creates a new TreeMap using a custom comparator
func NewWithComparator[K key, V val](opts ...Option[K, V]) (*TreeMap[K, V], error) {
	t := &TreeMap[K, V]{}

	for _, opt := range opts {
		opt(t)
	}

	if t.comparator == nil {
		return nil, errors.New("must provide a valid comparator")
	}

	return t, nil
}

type Option[K key, V val] func(t *TreeMap[K, V])

func WithCompareFunc[K key, V val](compareFunc Comparator[K]) func(t *TreeMap[K, V]) {
	return func(t *TreeMap[K, V]) {
		t.comparator = compareFunc
	}
}

func WithComparer[K key, V val](comparer Comparer[K]) func(t *TreeMap[K, V]) {
	compareFunc := func(a, b K) int {
		return comparer.CompareTo(b)
	}
	return func(t *TreeMap[K, V]) {
		t.comparator = compareFunc
	}
}

// Comparator allows keys to be compared for searching.
// should return -1 if (a < b), 0 if (a == b), +1 if (a > b)
type Comparator[K any] func(a, b K) int

// Comparer can be implemented to compare the key to the target.
// should return -1 if (a < b), 0 if (a == b), +1 if (a > b)
type Comparer[K key] interface {
	CompareTo(b K) int
}

func defaultComparator[key constraints.Ordered](a, b key) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

// Insert a new element with a key and value or update value on an existing key.
func (t *TreeMap[K, V]) Insert(key K, val V) {
	t.root = t.insert(t.root, key, val)
	t.root.color = black
	t.length++
}

// insert will recursively traverse down the tree and insert new node at leaf or
// update the value if key exists, then fix by doing rotation or color flip
func (t *TreeMap[K, V]) insert(cur *node[K, V], key K, val V) *node[K, V] {
	if cur == nil {
		cur = newNode(key, val, red)
		return cur
	}

	c := t.comparator(key, cur.key)
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

// Length returns the number of elements in the tree map.
func (t *TreeMap[K, V]) Length() int {
	return t.length
}

// Search by key and returns value if found, or the zero value and false if not found
func (t *TreeMap[K, V]) Search(key K) (V, bool) {
	cur := t.root
	for cur != nil {
		c := t.comparator(key, cur.key)
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

// String prints the tree in a visual format row by row.
func (t *TreeMap[K, V]) String() string {
	d := 0
	list := map[int][]string{}
	traverseByDepth(t.root, d, list)

	sb := strings.Builder{}
	sb.WriteString("\n----treemap----\n")
	for i := 1; i <= len(list); i++ {
		sb.WriteString(fmt.Sprintf("[depth %d]:  ", i))
		sb.WriteString(fmt.Sprintf("%v\n", strings.Join(list[i], " | ")))
	}

	return sb.String()
}

// Iterator returns a new iterator and starts at the first element.
func (t *TreeMap[K, V]) Iterator() *Iterator[K, V] {
	it := &Iterator[K, V]{tree: t}

	it.Begin()

	return it
}

// Iterator traverses through the treemap in sorted order.
type Iterator[K key, V val] struct {
	tree    *TreeMap[K, V]
	current *node[K, V]
}

// Begin moves iterator in front of first element.
func (it *Iterator[K, V]) Begin() {
	cur := it.tree.root
	for cur.left != nil {
		cur = cur.left
	}

	it.current = cur
}

// Last moves iterator in front of the last element.
func (it *Iterator[K, V]) Last() {
	cur := it.tree.root
	for cur.right != nil {
		cur = cur.right
	}

	it.current = cur
}

// End moves iterator to behind the last element.
func (it *Iterator[K, V]) End() {
	it.current = nil
}

// Next does an in-order traversal through a binary search tree.
func (it *Iterator[K, V]) Next() bool {
	cur := it.current
	if cur == nil {
		it.Begin()
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
		it.Last()
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
