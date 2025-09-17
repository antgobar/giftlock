package server

import (
	"giftlock/internal/middleware"
	"log"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

type Router interface {
	RegisterRoutes(mux *http.ServeMux)
}

func NewServer(addr string, mw middleware.Middleware, routers ...Router) *Server {
	mux := http.NewServeMux()

	for _, r := range routers {
		r.RegisterRoutes(mux)
	}

	handler := mw(mux)

	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (s *Server) Run() {
	log.Printf("Starting on %s\n", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s", err)
	}
}
