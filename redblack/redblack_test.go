package redblack

import (
	"testing"
)

// WIP test for insert
func TestInsert_ThreeNodes(t *testing.T) {
	rb := NewRedBlackBST()
	rb.Insert(1, "a")
	rb.Insert(2, "b")
	rb.Insert(3, "c")

	assertEqual(t, "b", rb.root.value)
	assertEqual(t, "a", rb.root.left.value)
	assertEqual(t, "c", rb.root.right.value)
}

func assertEqual(t *testing.T, want, got interface{}) {
	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}
