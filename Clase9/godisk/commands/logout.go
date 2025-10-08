package commands

import (
	"fmt"
	"godisk/state"
)

// ExecuteLogout maneja el comando de cierre de sesión
func ExecuteLogout() {
	// Verificar si hay una sesión activa
	if !state.CurrentUser.IsLoggedIn {
		fmt.Println("Error: No hay ninguna sesión activa")
		return
	}

	// Cerrar sesión
	username := state.CurrentUser.Username
	state.CurrentUser.IsLoggedIn = false
	state.CurrentUser.Username = ""
	state.CurrentUser.MountedId = ""
	// --- INICIO DE LA CORRECCIÓN ---
	// Limpiar también los IDs de usuario y grupo
	state.CurrentUser.Id = 0
	state.CurrentUser.Group = 0
	// --- FIN DE LA CORRECCIÓN ---

	fmt.Printf("Sesión cerrada exitosamente para el usuario '%s'\n", username)
}
