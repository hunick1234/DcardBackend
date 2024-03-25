package main

import (
	_ "net/http/pprof"

	"github.com/hunick1234/DcardBackend/handlers"
	"github.com/hunick1234/DcardBackend/myhttp"
	"github.com/hunick1234/DcardBackend/server"
	"github.com/hunick1234/DcardBackend/storage/pool"
)

func main() {
	httpServer := myhttp.DebugServer()

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
