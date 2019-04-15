package server

import (
	"testing"
)

func TestCreateFilter_add_a_new_filter(t *testing.T) {
	s := NewServer()
	if s.Run() != true {
		t.Fatalf("implement")
	}

	_, err := s.CreateFilter("storage", 100)
	if err != nil {
		t.Fatalf("Not able to create")
	}
}

func TestFilters_AddTestRemoveTest(t *testing.T) {
	s := NewServer()
	s.CreateFilter("TheMostAmazingFilter", 100)
	v := []byte("foo")
	s.Add("TheMostAmazingFilter", &v)
	if s.Test("TheMostAmazingFilter", &v) != true {
		t.Fatalf("Membership is not working")
	}
	s.Remove("TheMostAmazingFilter", &v)
	if s.Test("TheMostAmazingFilter", &v) != false {
		t.Fatalf("Membership removal is not working")
	}
}
