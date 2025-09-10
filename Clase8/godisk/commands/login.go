package commands

import (
	"fmt"
	"godisk/fs"
	"godisk/state"
	"os"
	"strconv"
	"strings"
)

// ExecuteLogin maneja el comando de inicio de sesión
func ExecuteLogin(user, pass, id string) {
	// Verificar si ya hay una sesión activa
	if state.CurrentUser.IsLoggedIn {
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

	// --- INICIO DE LA CORRECCIÓN ---
	// Ahora authenticateUser devuelve uid, gid y si fue exitoso.
	uid, gid, ok := authenticateUser(partition, user, pass)
	if ok {
		// Iniciar sesión exitosa, guardando TODOS los datos.
		state.CurrentUser.IsLoggedIn = true
		state.CurrentUser.Username = user
		state.CurrentUser.Id = uid
		state.CurrentUser.Group = gid
		state.CurrentUser.MountedId = id
		fmt.Printf("Sesión iniciada exitosamente para el usuario '%s' en la partición '%s'\n", user, id)
	} else {
		fmt.Println("Error: Usuario no encontrado o contraseña incorrecta")
	}
	// --- FIN DE LA CORRECCIÓN ---
}

// authenticateUser verifica las credenciales del usuario en users.txt
func authenticateUser(partition state.MountedPartition, username, password string) (int32, int32, bool) {
	// Caso especial para el superusuario root
	if username == "root" && password == "123" {
		return 1, 1, true
	}

	diskFile, err := os.Open(partition.Path)
	if err != nil {
		fmt.Printf("Error al abrir el disco: %v\n", err)
		return -1, -1, false
	}
	defer diskFile.Close()

	superblock, err := fs.ReadSuperblock(diskFile, int64(partition.Start))
	if err != nil {
		fmt.Printf("Error al leer el superbloque: %v\n", err)
		return -1, -1, false
	}

	usersInode, _, err := fs.FindInodeByPath(diskFile, superblock, "/users.txt")
	if err != nil {
		fmt.Printf("Error al buscar users.txt: %v\n", err)
		return -1, -1, false
	}

	content, err := fs.ReadFileContent(diskFile, superblock, usersInode)
	if err != nil {
		fmt.Printf("Error al leer users.txt: %v\n", err)
		return -1, -1, false
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Split(line, ",")
		// Formato Usuario: 1,U,root,root,123 -> id, U, username, groupname, pass
		if len(parts) == 5 && parts[1] == "U" && parts[2] == username && parts[4] == password {
			// Encontramos al usuario. Ahora necesitamos su UID y GID.
			uid, err := strconv.ParseInt(parts[0], 10, 32)
			if err != nil {
				continue // Línea mal formada
			}
			// Buscamos el GID del grupo al que pertenece
			gid := findGroupId(lines, parts[3])
			return int32(uid), gid, true
		}
	}

	return -1, -1, false
}

// findGroupId es una función auxiliar para encontrar el ID de un grupo por su nombre.
func findGroupId(lines []string, groupName string) int32 {
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, ",")
		// Formato Grupo: 1,G,root -> id, G, groupname
		if len(parts) == 3 && parts[1] == "G" && parts[2] == groupName {
			gid, err := strconv.ParseInt(parts[0], 10, 32)
			if err == nil {
				return int32(gid)
			}
		}
	}
	return -1 // Grupo no encontrado
}
