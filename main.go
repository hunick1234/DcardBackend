package main

import (
	"fmt"
	"net/http"

	"github.com/hunick1234/DcardBackend/Myhttp"
	v1 "github.com/hunick1234/DcardBackend/api/v1"
)

var server Myhttp.Server

func init() {
	server = Myhttp.Server{
		Server: &http.Server{
			Addr: ":8080",
		},
	}
}

func main() {

	server.SetRouter("/post/v1/ads", v1.CreatAD)
	server.SetRouter("/get/v1/ads", v1.GetAD)
	fmt.Println("Server is running on port " + server.Addr)
	err := server.Start()
	if err != nil {
		fmt.Println(err)
	}
}
