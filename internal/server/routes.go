package server

import "net/http"

func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./static")))

	// Endpoint para ejecutar comandos desde la terminal
	mux.HandleFunc("/api/command", s.HandleCommand)

	return mux
}
