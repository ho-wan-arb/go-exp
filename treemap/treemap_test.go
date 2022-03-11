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

	tm := New[int, string]()
	for _, kv := range m {
		tm.Insert(kv.k, kv.v)
	}

	for _, kv := range m {
		got, ok := tm.Search(kv.k)
		assertEqual(t, true, ok)
		assertEqual(t, kv.v, got, tm)
	}

	gotKey, ok := tm.Search(-1)
	assertEqual(t, false, ok)
	assertEqual(t, gotKey, "")

	gotLen := tm.Length()
	assertEqual(t, len(m), gotLen)
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
		{"b", "4"},
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

	tm := NewWithComparator[string, string](sortByStringLenFunc)
	for _, kv := range m {
		tm.Insert(kv.k, kv.v)
	}

	// should - iterator returns keys by order of string length
	it := tm.Iterator()
	it.Begin()
	assertEqual(t, "b", it.Key())
	assertEqual(t, "4", it.Value())

	ok := it.Next()
	assertEqual(t, true, ok)
	assertEqual(t, "aa", it.Key())

	ok = it.Next()
	assertEqual(t, true, ok)
	assertEqual(t, "ccc", it.Key())
}

func TestTreeMap_Iterate(t *testing.T) {
	t.Parallel()

	tm := New[int, string]()
	tm.Insert(4, "d")
	tm.Insert(3, "c")
	tm.Insert(1, "a")
	tm.Insert(2, "b")

	it := tm.Iterator()
	it.Begin()
	assertEqual(t, "a", it.Value(), tm)

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

	it.Last()
	assertEqual(t, "d", it.Value())

	it.End()
	assertEqual(t, 0, it.Key())
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
