package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router *chi.Mux
}

func New() *Server {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	fs := http.FileServer(http.Dir("web/www"))
	router.Handle("/www/*", http.StripPrefix("/www/", fs))

	router.Get("/", PageHandler)

	return &Server{
		router: router,
	}
}

func (server *Server) Run(port int) {
	log.Printf("Server listening on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), server.router)
}
