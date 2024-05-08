package myhttp

import (
	"context"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"strings"
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

type Response struct {
	Message   []byte
	StausCode int
}

type Router struct {
	path    string
	method  Method
	handler http.HandlerFunc
}

var server *Server

func init() {
	server = NewServer()
	server.setRouter("/debug/pprof/", http.MethodGet, pprof.Index)
	server.setRouter("/debug/pprof/cmdline", http.MethodGet, pprof.Cmdline)
	server.setRouter("/debug/pprof/profile", http.MethodGet, pprof.Profile)
	server.setRouter("/debug/pprof/symbol", http.MethodGet, pprof.Symbol)
	server.setRouter("/debug/pprof/trace", http.MethodGet, pprof.Trace)
}

func NewServer() *Server {
	return &Server{
		Server: &http.Server{
			Addr: ":8080",
		},
		MyMux: &MyMux{
			method: make(map[Method]map[Path]*Router, 10),
		},
	}
}

func DebugServer() *Server {
	return server
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
	if _, excite := m.method[Method(r.Method)]; !excite {
		return nil
	}
	if m.method[Method(r.Method)][Path(r.URL.Path)] == nil {
		log.Println("not found", r.Method, ":", r.URL.Path)
		return nil
	}

	log.Println("connect:", r.Method, r.URL.Path)
	return m.method[Method(r.Method)][Path(r.URL.Path)].handler
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	log.Println("-->", r.Method, r.URL.Path)
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	start := time.Now()
	if name, found := strings.CutPrefix(r.URL.Path, "/debug/pprof/"); found {
		if name != "" && !(name == "cmdline" || name == "profile" || name == "symbol" || name == "trace") {
			pprof.Index(w, r)
			return
		}
	}

	handler := useRouter(s.MyMux, r)
	if handler == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Println("use root request took", time.Since(start))

	start = time.Now()
	handler.ServeHTTP(w, r)
	log.Printf("request took %s", time.Since(start))
}

func (s *Server) Start() {
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
