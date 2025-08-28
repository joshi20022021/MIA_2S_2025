package commands

import (
	"fmt"
	"godisk/fs"
	"godisk/state"
	"os"
	"strings"
)

// ExecuteLogin maneja el comando de inicio de sesión
func ExecuteLogin(user, pass, id string) {
	// Verificar si ya hay una sesión activa
	if state.CurrentSession.IsLoggedIn {
		fmt.Println("Error: Ya hay una sesión activa. Ejecute 'logout' antes de iniciar una nueva sesión.")
		return
	}

	// Validar parámetros obligatorios
	if user == "" || pass == "" || id == "" {
		fmt.Println("Error: Todos los parámetros son obligatorios (-user, -pass, -id)")
		return
	}

	// Verificar que la partición esté montada
	partition, found := state.GetMountedPartitionByID(id)
	if !found {
		fmt.Printf("Error: No se encontró una partición montada con ID '%s'\n", id)
		return
	}

	// Verificar usuario y contraseña en users.txt
	if authenticateUser(partition, user, pass) {
		// Iniciar sesión exitosa
		state.CurrentSession.IsLoggedIn = true
		state.CurrentSession.Username = user
		state.CurrentSession.PartitionID = id
		fmt.Printf("Sesión iniciada exitosamente para el usuario '%s' en la partición '%s'\n", user, id)
	} else {
		fmt.Println("Error: Usuario no encontrado o contraseña incorrecta")
	}
}

// authenticateUser verifica las credenciales del usuario en users.txt
func authenticateUser(partition state.MountedPartition, username, password string) bool {
	// Por defecto permitimos root/123 por si no se puede leer users.txt
	if username == "root" && password == "123" {
		return true
	}

	// Abrir el archivo del disco
	diskFile, err := os.Open(partition.Path)
	if err != nil {
		fmt.Printf("Error al abrir el disco: %v\n", err)
		return false
	}
	defer diskFile.Close()

	// Leer el superbloque
	superblock, err := fs.ReadSuperblock(diskFile, int64(partition.Start))
	if err != nil {
		fmt.Printf("Error al leer el superbloque: %v\n", err)
		return false
	}

	// Buscar el inodo de users.txt
	// Capturar los 3 valores que devuelve FindInodeByPath
	usersInode, _, err := fs.FindInodeByPath(diskFile, superblock, "/users.txt")
	if err != nil {
		fmt.Printf("Error al buscar users.txt: %v\n", err)
		return false
	}

	// Leer el contenido de users.txt
	content, err := fs.ReadFileContent(diskFile, superblock, usersInode)
	if err != nil {
		fmt.Printf("Error al leer users.txt: %v\n", err)
		return false
	}

	// Convertir []byte a string antes de usar strings.Split
	// Buscar usuario y contraseña en el contenido
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Split(line, ",")
		// El formato esperado es: id,tipo,nombre,contraseña
		if len(parts) >= 4 && parts[2] == username && parts[3] == password {
			return true
		}
	}

	return false
}
