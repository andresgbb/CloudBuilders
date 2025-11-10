package main

import (
	"fmt"
	"log"
	"net/http"

	"cloudbuilders/internal/server"
)

func main() {
	fmt.Println("🚀 Iniciando servidor en http://localhost:8080")

	srv := server.New()
	if err := http.ListenAndServe(":8080", srv.Router()); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
