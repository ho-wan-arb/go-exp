package treemap

import (
	"testing"
)

func TestTreeMap_ValidBST(t *testing.T) {
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
	}

	tree := New[int, string]()
	for _, kv := range m {
		tree.Insert(kv.k, kv.v)
		tree.validateTree(t)
	}
}

// Check validity of red-black binary search tree.
func (rb *TreeMap[K, V]) validateTree(t *testing.T) {
	checkBST(t, rb)
	checkBalancedLinks(t, rb)
	checkSize(t, rb)
	check23Tree(t, rb)
}

func checkSize[K key, V val](t *testing.T, rb *TreeMap[K, V]) {
	heights := map[K]int{}
	if !isConsistentSize(rb.root, heights) {
		t.Errorf("not a balanced binary tree: heights: %v\n%v\n", heights, rb)
	}
}

// cache heights to avoid recomputing the same heights, only counting black links
func height[K key, V val](x *node[K, V], mk map[K]int) int {
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
func isConsistentSize[K key, V val](x *node[K, V], mk map[K]int) bool {
	if x == nil {
		return true
	}

	hl := height(x.left, mk)
	hr := height(x.right, mk)
	abs := abs(hl - hr)
	if abs > 1 {
		return false
	}
	return isConsistentSize(x.left, mk) && isConsistentSize(x.right, mk)
}

func checkBST[K key, V val](t *testing.T, rb *TreeMap[K, V]) {
	if !isBST(rb.root, nil, nil) {
		t.Errorf("not a valid Binary Search Tree\n%v\n", rb)
	}
}

// recursively check that every node is smaller or equal on left and larger or equal on right
func isBST[K key, V val](x *node[K, V], min, max *K) bool {
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

func checkBalancedLinks[K key, V val](t *testing.T, rb *TreeMap[K, V]) {
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
func isBalanced[K key, V val](x *node[K, V], black int) bool {
	if x == nil {
		return black == 0
	}

	if !x.isRed() {
		black--
	}

	return isBalanced(x.left, black) && isBalanced(x.right, black)
}

func check23Tree[K key, V val](t *testing.T, rb *TreeMap[K, V]) {
	if !is23Tree(rb.root) {
		t.Errorf("not a valid 23 Tree\n%v\n", rb)
	}
}

// cannot have red right link, or 2 left red links in a row
func is23Tree[K key, V val](x *node[K, V]) bool {
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

// numeric helpers
func max(source, target int) int {
	if source > target {
		return source
	}
	return target
}

func abs(num int) int {
	if num >= 0 {
		return num
	}
	return -num
}
