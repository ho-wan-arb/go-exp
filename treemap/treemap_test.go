package treemap

import (
	"fmt"
	"testing"
)

func TestTreeMap_InsertAndSearch(t *testing.T) {
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

	tr := New[int, string]()
	for _, kv := range m {
		tr.Insert(kv.k, kv.v)
	}

	for _, kv := range m {
		got, ok := tr.Search(kv.k)
		assertEqual(t, kv.v, got, tr)
		assertEqual(t, true, ok, tr)
	}

	got := tr.Length()
	assertEqual(t, len(m), got, tr)
}

func TestTreeMap_CustomComparator(t *testing.T) {
	t.Parallel()

	m := []struct {
		k string
		v string
	}{
		{"aa", "2"},
		{"b", "1"},
		{"ccc", "3"},
	}

	sortByStringLenFunc := func(a, b string) int {
		switch {
		case len(a) < len(b):
			return -1
		case len(a) > len(b):
			return 1
		default:
			return 0
		}
	}
	tr := NewWithComparator[string, string](sortByStringLenFunc)
	for _, kv := range m {
		tr.Insert(kv.k, kv.v)
	}

	// should - iterator returns keys by order of string length
	it := tr.Begin()
	assertEqual(t, "b", it.Key())

	ok := it.Next()
	assertEqual(t, true, ok)
	assertEqual(t, "aa", it.Key())

	ok = it.Next()
	assertEqual(t, true, ok)
	assertEqual(t, "ccc", it.Key())
}

func TestTreeMap_Iterate(t *testing.T) {
	t.Parallel()

	tr := New[int, string]()
	tr.Insert(4, "d")
	tr.Insert(3, "c")
	tr.Insert(1, "a")
	tr.Insert(2, "b")

	it := tr.Begin()
	assertEqual(t, "a", it.Value(), tr)

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

	it = tr.Last()
	assertEqual(t, "d", it.Value())

	it = tr.End()
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
