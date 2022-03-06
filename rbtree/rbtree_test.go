package rbtree

import (
	"testing"

	"golang.org/x/exp/constraints"
)

func TestRedBlackBST_Insert(t *testing.T) {
	t.Parallel()

	m := []struct {
		k int
		v string
	}{
		{1, "a"},
		{4, "d"},
		{3, "c"},
		{2, "b"},
	}

	tree := New[int, string]()
	for _, kv := range m {
		tree.Insert(kv.k, kv.v)
		tree.validateTree(t)
	}
}

func TestRedBlackBST_Search(t *testing.T) {
	t.Parallel()

	m := []struct {
		k int
		v string
	}{
		{4, "d"},
		{3, "c"},
		{1, "a"},
		{2, "b"},
	}

	tree := New[int, string]()
	for _, kv := range m {
		tree.Insert(kv.k, kv.v)
	}

	for _, kv := range m {
		got := tree.Search(kv.k)
		assertEqual(t, kv.v, got)
	}
}

// Check validity of red-black binary search tree
func (rb *RedBlackBST[K, V]) validateTree(t *testing.T) {
	checkBST(t, rb)
	checkBalancedLinks(t, rb)
	checkSize(t, rb)
	check23Tree(t, rb)
}

func checkSize[K Key, V Value](t *testing.T, rb *RedBlackBST[K, V]) {
	heights := map[K]int{}
	if !isConsistentSize(rb.root, heights) {
		t.Errorf("not a balanced binary tree")
	}
}

// cache heights to avoid recomputing the same heights
func height[K Key, V Value](x *Node[K, V], mk map[K]int) int {
	if x == nil {
		return 0
	}

	mh, ok := mk[x.key]
	if ok {
		return mh
	}

	hl := height(x.left, mk)
	hr := height(x.right, mk)
	h := Max(hl, hr) + 1

	mk[x.key] = h
	return h
}

func Max[C constraints.Ordered](source, target C) C {
	if source > target {
		return source
	}
	return target
}

// recursively check that max height of left subtree is at most 1 different from height of right
func isConsistentSize[K Key, V Value](x *Node[K, V], mk map[K]int) bool {
	if x == nil {
		return true
	}

	hl := height(x.left, mk)
	hr := height(x.right, mk)
	abs := Abs(hl - hr)
	if abs > 1 {
		return false
	}
	return isConsistentSize(x.left, mk) && isConsistentSize(x.right, mk)
}

func Abs[N constraints.Float | constraints.Integer](num N) N {
	if num >= 0 {
		return num
	}
	return -num
}

func checkBST[K Key, V Value](t *testing.T, rb *RedBlackBST[K, V]) {
	if !isBST(rb.root, nil, nil) {
		t.Errorf("not a valid Binary Search Tree")
	}
}

// recursively check that every node is smaller or equal on left and larger or equal on right
func isBST[K Key, V Value](x *Node[K, V], min, max *K) bool {
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

func checkBalancedLinks[K Key, V Value](t *testing.T, rb *RedBlackBST[K, V]) {
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
func isBalanced[K Key, V Value](x *Node[K, V], black int) bool {
	if x == nil {
		return black == 0
	}

	if !x.IsRed() {
		black--
	}

	return isBalanced(x.left, black) && isBalanced(x.right, black)
}

func check23Tree[K Key, V Value](t *testing.T, rb *RedBlackBST[K, V]) {
	if !is23Tree(rb.root) {
		t.Errorf("not a valid 23 Tree")
	}
}

// cannot have red right link, or 2 left red links in a row
func is23Tree[K Key, V Value](x *Node[K, V]) bool {
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

func assertEqual(t *testing.T, want, got interface{}) {
	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}
