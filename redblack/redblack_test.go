package redblack

import (
	"math"
	"testing"
)

// WIP test for insert
func TestInsert_ThreeNodes(t *testing.T) {
	rb := NewRedBlackBST()
	rb.Insert(1, "a")
	rb.Insert(2, "b")
	rb.Insert(3, "c")
	rb.Insert(4, "d")

	assertEqual(t, "b", rb.root.value)
	assertEqual(t, "a", rb.root.left.value)
	assertEqual(t, "d", rb.root.right.value)
	assertEqual(t, "c", rb.root.right.left.value)

	checkBST(t, rb)
	checkBalancedLinks(t, rb)
	checkSize(t, rb)
	check23Tree(t, rb)
}

func assertEqual(t *testing.T, want, got interface{}) {
	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}

// Check validity of red-black binary search tree

func checkSize(t *testing.T, rb *RedBlackBST) {
	heights := map[Key]int{}
	if !isConsistentSize(rb.root, heights) {
		t.Errorf("not a balanced binary tree")
	}
}

// cache heights to avoid recomputing the same heights
func height(x *Node, mk map[Key]int) int {
	if x == nil {
		return 0
	}

	mh, ok := mk[x.key]
	if ok {
		return mh
	}

	hl := height(x.left, mk)
	hr := height(x.right, mk)
	h := int(math.Max(float64(hl), float64(hr))) + 1

	mk[x.key] = h
	return h
}

// recursively check that max height of left subtree is at most 1 different from height of right
func isConsistentSize(x *Node, mk map[Key]int) bool {
	if x == nil {
		return true
	}

	hl := height(x.left, mk)
	hr := height(x.right, mk)
	abs := math.Abs(float64(hl) - float64(hr))
	if abs > 1 {
		return false
	}
	return isConsistentSize(x.left, mk) && isConsistentSize(x.right, mk)
}

func checkBST(t *testing.T, rb *RedBlackBST) {
	if !isBST(rb.root, nil, nil) {
		t.Errorf("not a valid Binary Search Tree")
	}
}

// recursively check that every node is smaller or equal on left and larger or equal on right
func isBST(x *Node, min, max *Key) bool {
	if x == nil {
		return true
	}

	if min != nil && CompareTo(x.key, *min) <= 0 {
		return false
	}
	if max != nil && CompareTo(x.key, *max) >= 0 {
		return false
	}

	return isBST(x.left, min, &x.key) && isBST(x.right, &x.key, max)
}

func checkBalancedLinks(t *testing.T, rb *RedBlackBST) {
	// count black links from root to left most leaf
	black := 0
	x := rb.root

	for x != nil {
		if !x.IsRed() {
			black++
		}
		x = x.left
	}

	if !isBalanced(rb.root, black) {
		t.Errorf("tree is not balanced: want depth of %v", black)
	}
}

// recursively check that every leaf has the same count of black links
func isBalanced(x *Node, black int) bool {
	if x == nil {
		return black == 0
	}

	if !x.IsRed() {
		black--
	}

	return isBalanced(x.left, black) && isBalanced(x.right, black)
}

func check23Tree(t *testing.T, rb *RedBlackBST) {
	if !is23Tree(rb.root) {
		t.Errorf("not a valid 23 Tree")
	}
}

// cannot have red right link, or 2 left red links in a row
func is23Tree(x *Node) bool {
	if x == nil {
		return true
	}

	if x.right.IsRed() {
		return false
	}

	if x.left.IsRed() && x.left.left.IsRed() {
		return false
	}

	return is23Tree(x.left) && is23Tree(x.right)
}
