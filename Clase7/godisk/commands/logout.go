package commands

import (
	"fmt"
	"godisk/state"
)

// ExecuteLogout maneja el comando de cierre de sesión
func ExecuteLogout() {
	// Verificar si hay una sesión activa
	if !state.CurrentSession.IsLoggedIn {
		fmt.Println("Error: No hay ninguna sesión activa")
		return
	}

	// Cerrar sesión
	username := state.CurrentSession.Username
	state.CurrentSession.IsLoggedIn = false
	state.CurrentSession.Username = ""
	state.CurrentSession.PartitionID = ""

	fmt.Printf("Sesión cerrada exitosamente para el usuario '%s'\n", username)
}
