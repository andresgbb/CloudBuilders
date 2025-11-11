package server

// Importamos solo lo necesario
// (ya no necesitamos fmt ni net/http en este archivo)
 
// Server representa nuestro servidor HTTP.
type Server struct{}

// New devuelve una nueva instancia del servidor.
func New() *Server {
	return &Server{}
}