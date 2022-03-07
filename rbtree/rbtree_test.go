package rbtree

import (
	"fmt"
	"testing"

	"golang.org/x/exp/constraints"
)

func TestBSTree_Insert(t *testing.T) {
	t.Parallel()

	m := []struct {
		k int
		v string
	}{
		{1, "a"},
		{4, "d"},
		{3, "c"},
		{2, "b"},
		{5, "e"},
		{6, "f"},
		{9, "i"},
		{8, "h"},
		{7, "g"},
	}

	tree := New[int, string]()
	for _, kv := range m {
		tree.Insert(kv.k, kv.v)
		tree.validateTree(t)
	}
	fmt.Print(tree)
}

func TestBSTree_Search(t *testing.T) {
	t.Parallel()

	m := []struct {
		k string
		v int
	}{
		{"d", 4},
		{"c", 3},
		{"a", 1},
		{"b", 2},
	}

	tree := New[string, int]()
	for _, kv := range m {
		tree.Insert(kv.k, kv.v)
	}

	for _, kv := range m {
		got, ok := tree.Search(kv.k)
		assertEqual(t, kv.v, got)
		assertEqual(t, true, ok)
	}
	fmt.Print(tree)
}

func TestBSTree_Iterate(t *testing.T) {
	t.Parallel()

	tree := New[int, string]()
	tree.Insert(4, "d")
	tree.Insert(3, "c")
	tree.Insert(1, "a")
	tree.Insert(2, "b")

	fmt.Print(tree)

	it := tree.Begin()
	assertEqual(t, "a", it.Value())

	// in-order traversal
	ok := it.Next()
	assertEqual(t, true, ok)
	assertEqual(t, "b", it.Value())
	ok = it.Next()
	assertEqual(t, true, ok)
	assertEqual(t, "c", it.Value())
	ok = it.Next()
	assertEqual(t, true, ok)
	assertEqual(t, "d", it.Value())
	ok = it.Next()
	assertEqual(t, false, ok)
	assertEqual(t, "", it.Value())

	ok = it.Next()
	assertEqual(t, true, ok)
	// default to zero value if at end
	assertEqual(t, "a", it.Value())

	it = tree.Last()
	assertEqual(t, "d", it.Value())

	it = tree.End()
	assertEqual(t, "", it.Value())

	// in-order traversal in revesse
	ok = it.Prev()
	assertEqual(t, true, ok)
	assertEqual(t, "d", it.Value())
	ok = it.Prev()
	assertEqual(t, true, ok)
	assertEqual(t, "c", it.Value())
	ok = it.Prev()
	assertEqual(t, true, ok)
	assertEqual(t, "b", it.Value())
	ok = it.Prev()
	assertEqual(t, true, ok)
	assertEqual(t, "a", it.Value())
	ok = it.Prev()
	assertEqual(t, false, ok)
	assertEqual(t, "", it.Value())
}

// Check validity of red-black binary search tree
func (rb *RBTree[K, V]) validateTree(t *testing.T) {
	checkBST(t, rb)
	checkBalancedLinks(t, rb)
	checkSize(t, rb)
	check23Tree(t, rb)
}

func checkSize[K Key, V Value](t *testing.T, rb *RBTree[K, V]) {
	heights := map[K]int{}
	if !isConsistentSize(rb.root, heights) {
		t.Errorf("not a balanced binary tree: heights: %v\n%v\n", heights, rb)
	}
}

// cache heights to avoid recomputing the same heights, only counting black links
func height[K Key, V Value](x *RBNode[K, V], mk map[K]int) int {
	if x == nil {
		return 0
	}

	mh, ok := mk[x.key]
	if ok {
		return mh
	}

	hl := height(x.left, mk)
	hr := height(x.right, mk)
	h := max(hl, hr) + 1
	if x.left.isRed() {
		h--
	}

	mk[x.key] = h
	return h
}

// recursively check that max height of left subtree is at most 1 different from height of right
func isConsistentSize[K Key, V Value](x *RBNode[K, V], mk map[K]int) bool {
	if x == nil {
		return true
	}

	hl := height(x.left, mk)
	hr := height(x.right, mk)
	abs := abs(hl - hr)
	if abs > 1 {
		// for debugging
		fmt.Printf("key: %v, hl: %v, hr: %v\n", x.key, hl, hr)
		return false
	}
	return isConsistentSize(x.left, mk) && isConsistentSize(x.right, mk)
}

func checkBST[K Key, V Value](t *testing.T, rb *RBTree[K, V]) {
	if !isBST(rb.root, nil, nil) {
		t.Errorf("not a valid Binary Search Tree\n%v\n", rb)
	}
}

// recursively check that every node is smaller or equal on left and larger or equal on right
func isBST[K Key, V Value](x *RBNode[K, V], min, max *K) bool {
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

func checkBalancedLinks[K Key, V Value](t *testing.T, rb *RBTree[K, V]) {
	// count black links from root to left most leaf
	black := 0
	x := rb.root

	for x != nil {
		if !x.isRed() {
			black++
		}
		x = x.left
	}

	if !isBalanced(rb.root, black) {
		t.Errorf("tree is not balanced: want depth of %v\n%v\n", black, rb)
	}
}

// recursively check that every leaf has the same count of black links
func isBalanced[K Key, V Value](x *RBNode[K, V], black int) bool {
	if x == nil {
		return black == 0
	}

	if !x.isRed() {
		black--
	}

	return isBalanced(x.left, black) && isBalanced(x.right, black)
}

func check23Tree[K Key, V Value](t *testing.T, rb *RBTree[K, V]) {
	if !is23Tree(rb.root) {
		t.Errorf("not a valid 23 Tree\n%v\n", rb)
	}
}

// cannot have red right link, or 2 left red links in a row
func is23Tree[K Key, V Value](x *RBNode[K, V]) bool {
	if x == nil {
		return true
	}

	if x.right.isRed() {
		return false
	}

	if x.left.isRed() && x.left.left.isRed() {
		return false
	}

	return is23Tree(x.left) && is23Tree(x.right)
}

// assert helpers
func assertEqual(t *testing.T, want, got interface{}) {
	t.Helper()
	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}

// numeric helpers
func max[N constraints.Ordered](source, target N) N {
	if source > target {
		return source
	}
	return target
}

func abs[N constraints.Signed](num N) N {
	if num >= 0 {
		return num
	}
	return -num
}
