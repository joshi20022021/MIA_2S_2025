package commands

import (
	"encoding/binary"
	"fmt"
	"godisk/fs"
	"godisk/state"
	"godisk/structs"
	"os"
)

// ExecuteCat muestra el contenido de uno o múltiples archivos
// Parámetros:
// - filePaths: lista de rutas de archivos a mostrar
func ExecuteCat(filePaths []string) {
	// Verificar que hay un usuario logueado
	if !state.CurrentUser.IsLoggedIn {
		fmt.Println("Error: Debe estar logueado para usar este comando.")
		return
	}

	// Verificar que al menos se especificó un archivo
	if len(filePaths) == 0 {
		fmt.Println("Error: Debe especificar al menos un archivo.")
		return
	}

	// Obtener la partición montada del usuario actual
	mountedPartition, found := state.GetMountedPartitionByID(state.CurrentUser.MountedId)
	if !found {
		fmt.Printf("Error: No se encontró la partición montada con ID '%s'.\n", state.CurrentUser.MountedId)
		return
	}

	// Abrir el archivo del disco
	file, err := os.OpenFile(mountedPartition.Path, os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("Error al abrir el disco: %v\n", err)
		return
	}
	defer file.Close()

	// Leer el superbloque
	file.Seek(mountedPartition.Start, 0)
	var sb structs.Superblock
	if err := binary.Read(file, binary.BigEndian, &sb); err != nil {
		fmt.Printf("Error al leer el superbloque: %v\n", err)
		return
	}

	// Procesar cada archivo
	for i, filePath := range filePaths {
		if filePath == "" {
			continue
		}

		// Buscar el inodo del archivo
		inode, _, err := fs.FindInodeByPath(file, sb, filePath)
		if err != nil {
			fmt.Printf("Error: No se encontró el archivo '%s': %v\n", filePath, err)
			if i < len(filePaths)-1 {
				fmt.Println() // Separador entre archivos
			}
			continue
		}

		// Verificar que es un archivo (no un directorio)
		if inode.I_type != 1 {
			fmt.Printf("Error: '%s' no es un archivo.\n", filePath)
			if i < len(filePaths)-1 {
				fmt.Println() // Separador entre archivos
			}
			continue
		}

		// Verificar permisos de lectura para el usuario actual, root tiene acceso a todo

		// Leer el contenido del archivo
		content, err := fs.ReadFileContent(file, sb, inode)
		if err != nil {
			fmt.Printf("Error al leer el archivo '%s': %v\n", filePath, err)
			if i < len(filePaths)-1 {
				fmt.Println() // Separador entre archivos
			}
			continue
		}

		// Mostrar el contenido
		fmt.Print(string(content))

		// Agregar salto de línea entre archivos si hay más de uno
		if i < len(filePaths)-1 {
			fmt.Println()
		}
	}
}
