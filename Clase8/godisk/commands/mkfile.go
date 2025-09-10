package commands

import (
	"fmt"
	"godisk/fs"
	"godisk/state"
	"io/ioutil"
)

// ExecuteMkfile recibe los parámetros parseados para crear un archivo.
func ExecuteMkfile(path string, rFlag bool, size int, contPath string) {
	// 1. Validar parámetros
	if path == "" {
		fmt.Println("ERROR (mkfile): El parámetro -path es obligatorio.")
		return
	}
	if size < 0 {
		fmt.Println("ERROR (mkfile): El valor de -size no puede ser negativo.")
		return
	}

	// 2. Validar sesión y partición
	if !state.CurrentUser.IsLoggedIn {
		fmt.Println("ERROR (mkfile): No hay un usuario logueado para ejecutar el comando.")
		return
	}

	// Se unifica la obtención de la partición en una sola llamada.
	mountedPartition, ok := state.GetMountedPartitionByID(state.CurrentUser.MountedId)
	if !ok {
		fmt.Printf("ERROR (mkfile): La partición montada con ID '%s' no fue encontrada.\n", state.CurrentUser.MountedId)
		return
	}

	// 3. Preparar contenido
	var content []byte
	// El parámetro -cont tiene prioridad
	if contPath != "" {
		fileContent, err := ioutil.ReadFile(contPath)
		if err != nil {
			fmt.Printf("ERROR (mkfile): No se pudo leer el archivo de contenido en '%s'. %v\n", contPath, err)
			return
		}
		content = fileContent
		fmt.Printf("INFO: Se usará el contenido del archivo local: %s\n", contPath)
	} else if size > 0 {
		// Generar contenido con 0-9 si -size se especifica y -cont no.
		content = make([]byte, size)
		for i := 0; i < size; i++ {
			content[i] = byte('0' + (i % 10))
		}
		fmt.Printf("INFO: Se generará un archivo con contenido de tamaño %d bytes.\n", size)
	}

	// 4. Llamar a la lógica del sistema de archivos
	err := fs.CreateFile(mountedPartition, path, rFlag, content)
	if err != nil {
		fmt.Printf("ERROR (mkfile): %v\n", err)
	} else {
		fmt.Printf("INFO: Archivo '%s' creado exitosamente.\n", path)
	}
}
