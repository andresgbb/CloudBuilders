package server

import (
	"encoding/json"
	"net/http"
)

type CommandRequest struct {
	Command string `json:"command"`
}

type CommandResponse struct {
	Result string `json:"result"`
}

func (s *Server) HandleCommand(w http.ResponseWriter, r *http.Request) {
	var req CommandRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error leyendo el comando", http.StatusBadRequest)
		return
	}

	result := HandleTerminalCommand(req.Command)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CommandResponse{Result: result})
}
