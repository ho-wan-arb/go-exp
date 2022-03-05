package rbtree

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

type Node struct {
	key   Key
	value Value
	left  *Node
	right *Node
	color color
}

func newNode(key Key, val Value, clr color) *Node {
	return &Node{
		key:   key,
		value: val,
		color: clr,
	}
}

type RedBlackBST struct {
	root *Node
}

// New creates a new instance of a Left-Leaning Red-Black BST.
func New() *RedBlackBST {
	return &RedBlackBST{}
}

// CompareTo returns > 0 if source is greater than target
func CompareTo(source, target Key) int {
	if source > target {
		return 1
	}
	if source < target {
		return -1
	}

	return 0
}

// Insert a new element
func (t *RedBlackBST) Insert(key Key, val Value) {
	t.root = t.insert(t.root, key, val)
	t.root.color = COLOR_BLACK
}

// insert will recursively traverse down the tree and insert new node at leaf or
// update the value if key exists, then fix by doing rotation or color flip
func (t *RedBlackBST) insert(h *Node, key Key, val Value) *Node {
	if h == nil {
		h = newNode(key, val, COLOR_RED)
		return h
	}

	c := CompareTo(key, h.key)
	switch {
	case c < 0:
		h.left = t.insert(h.left, key, val)
	case c > 0:
		h.right = t.insert(h.right, key, val)
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

// Search by key and returns value
func (t *RedBlackBST) Search(key Key) Value {
	return search(t.root, key)
}

func search(x *Node, key Key) Value {
	for x != nil {
		c := CompareTo(key, x.key)
		if c < 0 {
			x = x.left
			continue
		}

		if c > 0 {
			x = x.right
			continue
		}

		return x.value
	}

	return nil
}

func (t *RedBlackBST) Delete() {
	// TODO
}

// utility functions on Node
func (h *Node) IsRed() bool {
	if h == nil {
		return false
	}
	return bool(h.color)
}

func (h *Node) rotateLeft() *Node {
	x := h.right
	h.right = x.left
	x.left = h
	x.color = x.left.color
	x.left.color = COLOR_RED
	return x
}

func (h *Node) rotateRight() *Node {
	x := h.left
	h.left = x.right
	x.right = h
	x.color = x.right.color
	x.right.color = COLOR_RED
	return x
}

func (h *Node) flipColors() {
	h.color = !h.color
	h.left.color = !h.left.color
	h.right.color = !h.right.color
}
