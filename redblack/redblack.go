package redblack

// An implmentation of the left-leaning red-black 2-3 binary search tree (LLRB BST).
//
// References:
//   https://sedgewick.io/wp-content/themes/sedgewick/papers/2008LLRB.pdf
//   https://algs4.cs.princeton.edu/code/edu/princeton/cs/algs4/RedBlackBST.java

const (
	COLOR_RED   color = true
	COLOR_BLACK color = false
)

type color bool

type (
	// TODO: Make generic and use comparator
	Key   int
	Value interface{}
)

type rbNode struct {
	key   Key
	value Value
	left  *rbNode
	right *rbNode
	color color
}

func newRBNode(key Key, val Value, clr color) *rbNode {
	return &rbNode{
		key:   key,
		value: val,
		color: clr,
	}
}

type RedBlackBST struct {
	root *rbNode
	size int
}

func NewRedBlackBST() *RedBlackBST {
	return &RedBlackBST{}
}

// CompareTo returns 1 if source is greater than target
func CompareTo(source, target Key) int {
	if source > target {
		return 1
	}
	if source < target {
		return -1
	}

	return 0
}

func (t *RedBlackBST) Insert(key Key, val Value) {
	t.root = t.insert(t.root, key, val)
	t.root.color = COLOR_BLACK
}

func (t *RedBlackBST) insert(h *rbNode, key Key, val Value) *rbNode {
	if h == nil {
		h = newRBNode(key, val, COLOR_RED)
		return h
	}

	// compare to key of node being inserted and traverse the tree based on result
	c := CompareTo(key, h.key)
	switch {
	case c < 0:
		h.left = t.insert(h.left, key, val)
	case c > 0:
		h.right = t.insert(h.right, key, val)
	default:
		// if key already exists, then just update the value
		h.value = val
	}

	// fix to ensure links lean left
	if h.right.isRed() && !h.left.isRed() {
		h = h.rotateLeft()
	}
	if h.left.isRed() && h.left.left.isRed() {
		h = h.rotateRight()
	}
	if h.left.isRed() && h.right.isRed() {
		h.flipColors()
	}

	return h
}

func (t *RedBlackBST) Search() interface{} {
	// TODO
	return nil
}

func (t *RedBlackBST) DeleteMin() {
	// TODO
}

func (t *RedBlackBST) Delete() {
	// TODO
}

// utility
func (h *rbNode) isRed() bool {
	if h == nil {
		return false
	}
	return bool(h.color)
}

/*
 [h]           [x]
a   x    ->  h    c
   b c      a b
*/
func (h *rbNode) rotateLeft() *rbNode {
	x := h.right
	h.right = x.left
	x.left = h
	x.color = x.left.color
	x.left.color = COLOR_RED
	return x
}

/*
   [h]         [x]
 x    c  ->   a   h
a b              b c
*/
func (h *rbNode) rotateRight() *rbNode {
	x := h.left
	h.left = x.right
	x.right = h
	x.color = x.right.color
	x.right.color = COLOR_RED
	return h
}

func (h *rbNode) flipColors() {
	h.color = !h.color
	h.left.color = !h.left.color
	h.right.color = !h.right.color
}
