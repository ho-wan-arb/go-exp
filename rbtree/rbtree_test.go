package rbtree

import (
	"fmt"
	"strings"
	"testing"
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
}

func TestRedBlackBST_Search(t *testing.T) {
	t.Parallel()

	m := []struct {
		k string
		v int
	}{
		{"d", 4},
		{"c", 3},
		{"aa", 1},
		{"ab", 2},
	}

	tree := New[string, int]()
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
		t.Errorf("not a balanced binary tree: heights: %v\n%v\n", heights, printByDepth(rb))
	}
}

// cache heights to avoid recomputing the same heights, only counting black links
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
	if x.left.IsRed() {
		h--
	}

	mk[x.key] = h
	return h
}

func Max(source, target int) int {
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
		// for debugging
		fmt.Printf("key: %v, hl: %v, hr: %v\n", x.key, hl, hr)
		return false
	}
	return isConsistentSize(x.left, mk) && isConsistentSize(x.right, mk)
}

func Abs(num int) int {
	if num >= 0 {
		return num
	}
	return -num
}

func checkBST[K Key, V Value](t *testing.T, rb *RedBlackBST[K, V]) {
	if !isBST(rb.root, nil, nil) {
		t.Errorf("not a valid Binary Search Tree\n%v\n", printByDepth(rb))
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
		t.Errorf("tree is not balanced: want depth of %v\n%v\n", black, printByDepth(rb))
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
		t.Errorf("not a valid 23 Tree\n%v\n", printByDepth(rb))
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

func printByDepth[K Key, V Value](rb *RedBlackBST[K, V]) string {
	d := 0
	list := map[int][]string{}
	traverseByDepth(rb.root, d, list)

	sb := strings.Builder{}
	for i := 1; i <= len(list); i++ {
		sb.WriteString(fmt.Sprintf("[depth %d]:  ", i))
		sb.WriteString(fmt.Sprintf("%v\n", strings.Join(list[i], " | ")))
	}
	sb.WriteString("----\n")

	return sb.String()
}

func traverseByDepth[K Key, V Value](x *Node[K, V], d int, list map[int][]string) {
	if x == nil {
		return
	}

	curKey := fmt.Sprintf("%v", x.key)

	if !x.IsRed() {
		d++
		list[d] = append(list[d], curKey)
	} else {
		// join 2 nodes: red link should lean left, so smaller number should always be in front
		list[d][len(list[d])-1] = fmt.Sprintf("(%v,%v)", curKey, list[d][len(list[d])-1])
	}

	traverseByDepth(x.left, d, list)
	traverseByDepth(x.right, d, list)
}
