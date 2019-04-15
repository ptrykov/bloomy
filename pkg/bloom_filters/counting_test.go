package filter

import (
	"testing"
)

func setupTest() *Counting {
	c := NewCounting(10)
	value := []byte("sebas 1")
	c.Add(&value)
	return c
}

func TestTest_return_true_for_members(t *testing.T) {
	c := setupTest()
	v := []byte("Pavel")
	c.Add(&v)
	if c.Test(&v) != true {
		t.Fatalf("The filter is not properly checking membership")
	}
}

func TestTest_return_false_for_non_members(t *testing.T) {
	c := setupTest()
	v := []byte("newvalue")
	if c.Test(&v) != false {
		t.Fatalf("Check false positive rates because this value should not be here")
	}
}

func TestRemove_return_false_if_not_a_member(t *testing.T) {
	c := setupTest()
	v := []byte("unseenvalue")
	if c.Remove(&v) != false {
		t.Fatalf("You removed a value that was not a member")
	}
}

func TestRemove_return_true_for_a_member(t *testing.T) {
	c := setupTest()
	v := []byte("insertme")
	c.Add(&v)
	if c.Remove(&v) != true {
		t.Fatalf("Filter didn't find the member to get removed")
	}

	if c.Test(&v) != false {
		t.Fatalf("Filter didn't remove the member")
	}
}
