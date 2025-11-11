package server

import "net/http"

// Router define las rutas de la aplicación.
func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()

	// Servimos los archivos estáticos (frontend)
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	// Endpoint para ejecutar comandos desde la terminal
	mux.HandleFunc("/api/command", s.HandleCommand)

	return mux
}

