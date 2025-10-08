package commands

import (
	"fmt"
	"godisk/state"   // Importamos el paquete de estado para acceder a la lista global.
	"godisk/structs" // Importamos las estructuras de MBR, EBR, etc.
	"godisk/utils"   // Importamos las herramientas para leer/escribir en el disco.
	"os"
	"strconv" // Para convertir números a texto (string).
	"strings"
)

// Estas son variables globales para llevar la cuenta de los IDs de montaje.
// Viven en memoria mientras el programa se ejecuta.
var diskLetters = make(map[string]rune)     // Mapa para asignar una letra a cada ruta de disco.
var nextLetter rune = 'A'                   // La siguiente letra disponible para un nuevo disco.
var partitionNumbers = make(map[string]int) // Mapa para llevar el número de la próxima partición por disco.

// ExecuteMount monta una partición en memoria.
func ExecuteMount(path, name string) {
	// --- 1. Verificar si la partición ya está montada ---
	// Recorre la lista global de particiones activas.
	for _, p := range state.GlobalMountedPartitions {
		// Si encuentra una partición con la misma ruta y nombre, ya está montada.
		if p.Path == path && p.Name == name {
			fmt.Printf("Error: la partición '%s' en el disco '%s' ya está montada.\n", name, path)
			return // Termina la función.
		}
	}

	// --- 2. Abrir y leer el disco para encontrar la partición ---
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("Error: no se pudo abrir el disco en '%s'.\n", path)
		return
	}
	defer file.Close() // Asegura que el archivo se cierre al final.

	// Lee el MBR del disco para poder buscar la partición.
	mbr, err := utils.ReadMBR(file)
	if err != nil {
		fmt.Printf("Error al leer el MBR: %v\n", err)
		return
	}

	// --- 3. Generar el ID único ---
	// Verifica si el disco ya tiene una letra asignada.
	if _, ok := diskLetters[path]; !ok {
		// Si es un disco nuevo, le asigna la siguiente letra disponible.
		diskLetters[path] = nextLetter
		nextLetter++ // Incrementa la letra para la próxima vez (A -> B -> C).
		// Inicia el contador de particiones para este nuevo disco en 1.
		partitionNumbers[path] = 1
	}
	// Obtiene la letra y el número de partición que le corresponde.
	letter := diskLetters[path]
	partNum := partitionNumbers[path]
	// Incrementa el número para la siguiente partición que se monte de ESTE MISMO disco.
	partitionNumbers[path]++

	// Construye el ID según la fórmula del enunciado.
	carnet := "202401234" // Carnet de ejemplo.
	// carnet[len(carnet)-2:] toma los últimos 2 dígitos.
	// strconv.Itoa convierte el número de partición a texto.
	id := carnet[len(carnet)-2:] + strconv.Itoa(partNum) + string(letter)

	// --- 4. Buscar y montar particiones primarias ---
	// Recorre las 4 particiones del MBR.
	for i := range mbr.Mbr_partitions {
		p := &mbr.Mbr_partitions[i] // Usamos un puntero para poder modificar la partición.
		// Compara el nombre, quitando los bytes nulos de relleno.
		if strings.Trim(string(p.Part_name[:]), "\x00") == name {
			// No se pueden montar particiones extendidas.
			if p.Part_type == 'E' {
				fmt.Println("Error: no se pueden montar particiones extendidas.")
				return
			}
			// --- Actualiza el disco ---
			p.Part_status = '1'
			p.Part_correlative = int64(partNum)
			copy(p.Part_id[:], id) // Copia el ID al campo de la partición.

			// --- Actualiza la memoria ---
			// Crea la "ficha" de la partición montada.
			newMount := state.MountedPartition{
				ID:      id,
				Path:    path,
				Name:    name,
				Status:  '1',
				Letter:  letter,
				PartNum: partNum,
				Size:    p.Part_s,     // <-- AÑADIR ESTA LÍNEA
				Start:   p.Part_start, // <-- AÑADIR ESTA LÍNEA
			}
			// La añade a la lista global.
			state.GlobalMountedPartitions = append(state.GlobalMountedPartitions, newMount)

			// Guarda los cambios (status, id, correlativo) en el archivo .mia.
			if err := utils.WriteMBR(file, &mbr); err != nil {
				fmt.Printf("Error al actualizar el MBR en el disco: %v\n", err)
				return
			}
			fmt.Printf("Partición primaria '%s' montada exitosamente con el ID: %s\n", name, id)
			return // Termina porque ya encontró y montó la partición.
		}
	}

	// --- 5. Si no es primaria, buscar y montar particiones lógicas ---
	var extendedPartition structs.Partition
	foundExtended := false
	// Busca si hay una partición extendida en el MBR.
	for i := range mbr.Mbr_partitions {
		if mbr.Mbr_partitions[i].Part_type == 'E' {
			extendedPartition = mbr.Mbr_partitions[i]
			foundExtended = true
			break
		}
	}

	// Si hay una partición extendida, busca la lógica dentro de ella.
	if foundExtended {
		// Lee el primer EBR al inicio de la partición extendida.
		currentEBR, err := utils.ReadEBR(file, extendedPartition.Part_start)
		if err != nil {
			fmt.Printf("Error al leer el primer EBR: %v\n", err)
			return
		}
		currentEBRAddress := extendedPartition.Part_start

		// Inicia el recorrido de la "lista enlazada" de EBRs.
		for {
			// Compara el nombre del EBR actual con el buscado.
			if strings.Trim(string(currentEBR.Part_name[:]), "\x00") == name {
				// --- Actualiza el disco ---
				currentEBR.Part_status = '1' // Solo actualizamos el status en el EBR.
				// --- Actualiza la memoria ---
				newMount := state.MountedPartition{
					ID:      id,
					Path:    path,
					Name:    name,
					Status:  '1',
					Letter:  letter,
					PartNum: partNum,
					Size:    currentEBR.Part_s,     // <-- AÑADIR ESTA LÍNEA
					Start:   currentEBR.Part_start, // <-- AÑADIR ESTA LÍNEA
				}
				state.GlobalMountedPartitions = append(state.GlobalMountedPartitions, newMount)

				// Guarda el EBR modificado en su lugar.
				if err := utils.WriteEBR(file, &currentEBR, currentEBRAddress); err != nil {
					fmt.Printf("Error al actualizar el EBR en el disco: %v\n", err)
					return
				}
				fmt.Printf("Partición lógica '%s' montada exitosamente con el ID: %s\n", name, id)
				return // Termina.
			}
			// Si no es el EBR que buscamos y no hay más en la cadena, termina el bucle.
			if currentEBR.Part_next == -1 {
				break
			}
			// Si hay más, actualiza la dirección y lee el siguiente EBR.
			currentEBRAddress = currentEBR.Part_next
			currentEBR, err = utils.ReadEBR(file, currentEBRAddress)
			if err != nil {
				fmt.Printf("Error al leer la cadena de EBRs: %v\n", err)
				return
			}
		}
	}

	// Si llega hasta aquí, es porque no encontró la partición en ningún lado.
	fmt.Printf("Error: no se encontró la partición con el nombre '%s'.\n", name)
}

// ExecuteMounted muestra todas las particiones montadas.
func ExecuteMounted() {
	// Revisa si la lista global está vacía.
	if len(state.GlobalMountedPartitions) == 0 {
		fmt.Println("No hay particiones montadas.")
		return
	}
	// Imprime un encabezado.
	fmt.Println("--- Particiones Montadas ---")
	// Recorre la lista e imprime los datos de cada partición montada.
	for _, p := range state.GlobalMountedPartitions {
		fmt.Printf("- ID: %s, Disco: %s, Partición: %s\n", p.ID, p.Path, p.Name)
	}
	fmt.Println("--------------------------")
}
