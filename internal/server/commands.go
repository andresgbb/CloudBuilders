package server

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// HandleTerminalCommand ejecutara la lógica de cada comando recibido
func HandleTerminalCommand(command string) string {
	switch {
	case command == "help":
		return `comandos disponibles:
		help	- Muestra esta ayuda
		date	- Muestra la fecha y hora actual
		os		- Muestra la información del SO
		whoami	- Muestra el usuario actual
		clear	- Limpia la pantalla
		About	- Información del proyecto
		echo X	- repite el texto que escribas despues de echo
		ls		- Lista archivos en el directorio actual (Demo)`

	case command == "date":
		return time.Now().Format("02/01/2006 15:04:05")

	case command == "os":
		return fmt.Sprintf("Sistema operativo: %s (%s)", runtime.GOOS, runtime.GOARCH)

	case command == "whoami":
		user := os.Getenv("USERNAME")
		if user == "" {
			user = os.Getenv("USER")
		}
		return fmt.Sprintf("Usuario actual: %s", user)

	case command == "about":
		return `CloudBuilders Terminal v1.0
		Proyecto educativo en Go + Docker + Kubernetes + AWS.
		Simula una terminal web interactiva.`

	//Comandos Dinamicos

	case strings.HasPrefix(command, "echo "):
		return strings.TrimPrefix(command, "echo ")

	case command == "ls":
		out, err := exec.Command("ls").Output()
		if err != nil {
			return fmt.Sprintf("Error ejecutando ls: %v", err)
		}
		return string(out)

	default:
		return fmt.Sprintf("Comando no reconocido: %s. Escribe 'help' para ver las opciones", command)
	}
}
