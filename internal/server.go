package server

import (
	"log"

	filter "github.com/ptrykov/bloomy/pkg"
	bf "github.com/ptrykov/bloomy/pkg/bloom_filters"
)

type Server struct {
	// config
	// map string -> filter
	config  *ServerConfig
	filters map[string]filter.Filter
}

func NewServer(cfg *ServerConfig) *Server {
	return &Server{
		config:  cfg,
		filters: make(map[string]filter.Filter),
	}
}

func (s *Server) Run() bool {
	log.Println("Listening on:", s.config.Port)
	return true
}

func (s *Server) CreateFilter(name string, size uint) (filter.Filter, error) {
	if _, ok := s.filters[name]; ok != true {
		s.filters[name] = bf.NewCounting(size)
	}
	return s.filters[name], nil
}

func (s *Server) loadFilters() {

}

func (s *Server) DeleteFilter(name string) {
	delete(s.filters, name)
}

func (s *Server) Add(name string, value *[]byte) {
	s.filters[name].Add(value)
}

func (s *Server) Test(name string, value *[]byte) bool {
	return s.filters[name].Test(value)
}

func (s *Server) Remove(name string, value *[]byte) bool {
	return s.filters[name].Remove(value)
}
