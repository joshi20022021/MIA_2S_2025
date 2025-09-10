package commands

import (
	"fmt"
	"os"
)

// ExecuteRmdisk contiene la lógica para eliminar un disco directamente.
func ExecuteRmdisk(path string) {
	// Intenta eliminar el archivo especificado en la ruta.
	err := os.Remove(path)
	if err != nil {
		// Verifica si el error es porque el archivo no existe.
		if os.IsNotExist(err) {
			fmt.Printf("Error: el archivo en la ruta '%s' no existe.\n", path)
		} else {
			// Informa de otros posibles errores (ej. falta de permisos).
			fmt.Printf("Error al eliminar el archivo: %v\n", err)
		}
		return
	}

	// Si no hubo errores, la eliminación fue exitosa.
	fmt.Printf("Disco en '%s' eliminado exitosamente.\n", path)
}
