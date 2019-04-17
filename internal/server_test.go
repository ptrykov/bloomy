package server

import (
	"testing"
)

func setupServer() *Server {
	cfg := &ServerConfig{Port: 333}
	return NewServer(cfg)
}

func TestCreateFilter_add_a_new_filter(t *testing.T) {
	s := setupServer()
	if s.Run() != true {
		t.Fatalf("implement")
	}

	_, err := s.CreateFilter("storage", 100)
	if err != nil {
		t.Fatalf("Not able to create")
	}
}

func TestCreateFilter_should_not_recreate_a_filter(t *testing.T) {
	s := setupServer()
	s.CreateFilter("one", 1)
	v := []byte("onetwothree")
	s.Add("one", &v)
	if s.Test("one", &v) != true {
		t.Fatalf("Broken membership")
	}
	s.CreateFilter("one", 1)
	if s.Test("one", &v) != true {
		t.Fatalf("The filter got resetted")
	}
}

func TestFilters_AddTestRemoveTest(t *testing.T) {
	s := setupServer()
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
