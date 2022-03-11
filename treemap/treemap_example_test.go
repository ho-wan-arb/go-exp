package treemap_test

import (
	"fmt"

	"github.com/ho-wan-arb/go-exp/treemap"
)

func ExampleNew() {
	t := treemap.New[int, string]()

	t.Insert(5, "apple")
	t.Insert(3, "banana")

	value, ok := t.Search(5)
	fmt.Println("found:", ok)
	fmt.Println("value:", value)

	// Output:
	// found: true
	// value: apple
}

func ExampleNewWithComparator() {
	sortByStringLenFunc := func(a, b string) int {
		switch {
		case len(a) == len(b):
			return 0
		case len(a) < len(b):
			return -1
		default:
			return 1
		}
	}

	t := treemap.NewWithComparator[string, string](sortByStringLenFunc)

	t.Insert("aaa", "apple")
	t.Insert("b", "banana")

	it := t.Iterator()
	it.Begin()

	// items are returned in sorted order of key length
	fmt.Println("item 1:", it.Value())

	fmt.Println("has next:", it.Next())
	fmt.Println("item 2:", it.Value())

	fmt.Println("has next:", it.Next())
	fmt.Println("item 3:", it.Value())

	// Output:
	// item 1: banana
	// has next: true
	// item 2: apple
	// has next: false
	// item 3:
}

func ExamplePrint() {
	type item struct{}
	t := treemap.New[int, *item]()

	t.Insert(3, nil)
	t.Insert(4, nil)
	t.Insert(2, &item{})
	t.Insert(1, nil)

	fmt.Println(t)

	// Output:
	// ----treemap----
	// [depth 1]:  3
	// [depth 2]:  (1,2) | 4
}
