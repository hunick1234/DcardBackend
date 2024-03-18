package server

import (
	"github.com/hunick1234/DcardBackend/myhttp"
	"github.com/hunick1234/DcardBackend/storage/pool"
)

type Server struct {
	Pool *pool.Pool
	HTTP *myhttp.Server
}

func NewAdServer(s pool.Pool, http myhttp.Server) Server {
	return Server{}
}

func (s *Server) Run() {
	s.HTTP.Start()
}
