package server

import (
	"fmt"
	"net/http"
)

type Server struct{}

func New() *Server {
	return &Server{}
}

func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "👋 Bienvenido a CloudBuilders Terminal!")
	})

	return mux
}
