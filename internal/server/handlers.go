package server

import (
	"encoding/json"
	"net/http"
)

// Estructura que representa la solicitud del cliente (lo que escribe el usuario en la terminal)
type CommandRequest struct {
	Command string `json:"command"`
}

// Estructura para enviar la respuesta al cliente (el resultado del comando)
type CommandResponse struct {
	Result string `json:"result"`
}

// HandleCommand maneja la ruta"/api/command"
// Recibe un comando JSON y devuelve una respuesta JSON con el resultado
func (s *Server) HandleCommand(w http.ResponseWriter, r *http.Request) {
	var req CommandRequest

	// Decodificamos el cuerpo JSON de la petición
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error leyendo el comando", http.StatusBadRequest)
		return
	}

	// Ejecutamos el comando (se maneja en commands.go)
	result := HandleTerminalCommand(req.Command)

	// Creamos la respuesta JSON
	// resp := CommandResponse{Result: result}

	// Codificamos y enviamos la respuesta al cliente
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CommandResponse{Result: result})
}
