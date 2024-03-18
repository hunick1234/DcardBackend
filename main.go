package main

import (
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
}
