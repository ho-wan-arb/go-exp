package treemap

import (
	"fmt"
	"testing"
)

func TestTreeMap_Search(t *testing.T) {
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
	}

	for _, kv := range m {
		got, ok := tree.Search(kv.k)
		assertEqual(t, kv.v, got, tree)
		assertEqual(t, true, ok, tree)
	}
}

func TestTreeMap_Iterate(t *testing.T) {
	t.Parallel()

	tree := New[int, string]()
	tree.Insert(4, "d")
	tree.Insert(3, "c")
	tree.Insert(1, "a")
	tree.Insert(2, "b")

	it := tree.Begin()
	assertEqual(t, "a", it.Value(), tree)

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

// assert helpers
func assertEqual(t *testing.T, want, got any, msgAndArgs ...interface{}) {
	t.Helper()
	if want != got {
		t.Errorf(fmt.Sprintf("want %v, got %v", want, got), msgAndArgs...)
	}
}
