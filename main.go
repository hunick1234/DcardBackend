package main

import (
<<<<<<< HEAD
	_ "net/http/pprof"

	"github.com/hunick1234/DcardBackend/Myhttp"
	v1 "github.com/hunick1234/DcardBackend/api/v1"
)

var server *Myhttp.Server

func init() {
	server = Myhttp.NewServer()
	server.Addr = ":8080"

}

func main() {

	server.Get("/api/v1/ads", v1.GetAD)
	server.Post("/api/v1/ads", v1.CreatAD)
	server.Start()
=======
	"github.com/hunick1234/DcardBackend/handlers"
	"github.com/hunick1234/DcardBackend/myhttp"
	"github.com/hunick1234/DcardBackend/server"
	"github.com/hunick1234/DcardBackend/storage/pool"
)

func main() {
	httpServer := myhttp.NewServer()

	connPool := pool.NewPool()
	defer connPool.ClosePool()
	// mockService:=service.NewMockService()
	server := server.Server{
		Pool: connPool,
		HTTP: httpServer,
	}
	handlers.StartAPIControll(&server)

	server.Run()
>>>>>>> develop
}
