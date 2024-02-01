package Myhttp

import (
	"net/http"
)

type Server struct {
	*http.Server
}

type HttpMethod interface {
	Get(path string, router http.HandlerFunc)
	Post(path string, router http.HandlerFunc)
}

func (s *Server) SetRouter(path string, router http.HandlerFunc) {

	http.HandleFunc(path, router)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) Start() error {
	return http.ListenAndServe(s.Addr, nil)
}

func (s *Server) Stop() error {
	return s.Close()
}

func (s *Server) Method() HttpMethod {
	return s
}

func (s *Server) Get(path string, router http.HandlerFunc) {

	http.HandleFunc(path, router)
}

func (s *Server) Post(path string, router http.HandlerFunc) {

	http.HandleFunc(path, router)
}
