package commands

import (
	"fmt"
	"godisk/fs"
	"godisk/state"
)

// ExecuteMkdir crea un directorio - path: ruta del directorio a crear - recursive: flag -p para crear directorios padre si no existen
func ExecuteMkdir(path string, recursive bool) {
	// Verificar que hay un usuario logueado
	if !state.CurrentUser.IsLoggedIn {
		fmt.Println("Error: Debe estar logueado para usar este comando.")
		return
	}

	// Verificar que el parámetro es válido
	if path == "" {
		fmt.Println("Error: Debe especificar la ruta del directorio.")
		return
	}

	// Obtener la partición montada del usuario actual
	mountedPartition, found := state.GetMountedPartitionByID(state.CurrentUser.MountedId)
	if !found {
		fmt.Printf("Error: No se encontró la partición montada con ID '%s'.\n", state.CurrentUser.MountedId)
		return
	}

	// Crear el directorio usando fs.CreateDirectory
	err := fs.CreateDirectory(mountedPartition, path, recursive)
	if err != nil {
		fmt.Printf("Error al crear el directorio: %v\n", err)
		return
	}

	fmt.Printf("Directorio '%s' creado exitosamente.\n", path)
}
