package Myhttp

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	*http.Server
	*MyMux
}
type Method string
type Path string
type MyMux struct {
	method map[Method]map[Path]*Router
}

type Router struct {
	path    string
	method  Method
	handler http.HandlerFunc
}

func NewServer() Server {
	return Server{
		Server: &http.Server{
			Addr: ":8080",
		},
		MyMux: &MyMux{
			method: make(map[Method]map[Path]*Router, 10),
		},
	}
}

// todo check path first byte is '/'
func (s *Server) setRouter(path string, method Method, handler http.HandlerFunc) {
	if s.MyMux.method[method] == nil {
		s.MyMux.method[method] = make(map[Path]*Router, 10)
	}

	s.method[method][Path(path)] = &Router{
		path:    path,
		method:  method,
		handler: handler,
	}
}

func useRouter(m *MyMux, r *http.Request) http.Handler {
	log.Println("connect:", r.Method, r.URL.Path)

	if m.method[Method(r.Method)] == nil {
		return nil
	}

	if m.method[Method(r.Method)][Path(r.URL.Path)] == nil {
		log.Println("not found", r.Method, ":", r.URL.Path)
		return nil
	}

	return m.method[Method(r.Method)][Path(r.URL.Path)].handler
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	handler := useRouter(s.MyMux, r)
	if handler == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	handler.ServeHTTP(w, r)
}

func (s *Server)  Start() {
	go func() {
		if err := http.ListenAndServe(s.Addr, s); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	log.Printf("Server is running on port %s", s.Addr)
	waitSignal()
	greatFullShutdown(s)
}

func (s *Server) Stop() error {
	return s.Close()
}

func (s *Server) Get(path string, handler http.HandlerFunc) {
	log.Println("Get:", path)
	s.setRouter(path, http.MethodGet, handler)
}

func (s *Server) Post(path string, handler http.HandlerFunc) {
	log.Println("Post:", path)
	s.setRouter(path, http.MethodPost, handler)
}

func waitSignal() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server is shutting down...")
}

func greatFullShutdown(server *Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
