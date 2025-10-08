package commands

import (
	"encoding/binary"
	"fmt"
	"godisk/state"
	"godisk/structs"
	"godisk/utils"
	"os"
	"strings"
)

// ExecuteJournal muestra las entradas del journal para una partición específica
func ExecuteJournal(id string) {
	// Validar que se proporcionó el ID
	if id == "" {
		fmt.Println("ERROR (journal): El parámetro -id es obligatorio.")
		return
	}

	// Buscar la partición montada
	mountedPartition, ok := state.GetMountedPartitionByID(id)
	if !ok {
		fmt.Printf("Error: No se encontró la partición montada con ID '%s'.\n", id)
		return
	}

	// Abrir el archivo del disco
	file, err := os.OpenFile(mountedPartition.Path, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Printf("Error: No se pudo abrir el disco: %v\n", err)
		return
	}
	defer file.Close()

	// Leer el superbloque usando BigEndian
	var sb structs.Superblock
	_, err = file.Seek(mountedPartition.Start, 0)
	if err != nil {
		fmt.Printf("Error: No se pudo posicionar en el superbloque: %v\n", err)
		return
	}

	err = binary.Read(file, binary.BigEndian, &sb)
	if err != nil {
		fmt.Printf("Error: No se pudo leer el superbloque: %v\n", err)
		return
	}

	// Leer las entradas del journal
	entries, err := utils.ReadJournal(file, &sb)
	if err != nil {
		fmt.Printf("Error: No se pudo leer el journal: %v\n", err)
		return
	}

	// Mostrar las entradas en formato tabla
	ShowJournalReport(entries)
}

// ShowJournalReport muestra las entradas del journal en formato de tabla
func ShowJournalReport(entries []structs.JournalEntry) {
	fmt.Println("═══════════════════════════════════════════════════════════════════")
	fmt.Println("                            JOURNAL DE OPERACIONES                 ")
	fmt.Println("═══════════════════════════════════════════════════════════════════")
	fmt.Printf("%-15s %-40s %-30s %-20s\n", "OPERACION", "PATH", "CONTENIDO", "FECHA")
	fmt.Println("───────────────────────────────────────────────────────────────────")

	if len(entries) == 0 {
		fmt.Println("No hay entradas en el journal.")
		fmt.Println("═══════════════════════════════════════════════════════════════════")
		return
	}

	for _, entry := range entries {
		// Convertir arrays de bytes a strings y limpiar caracteres nulos
		operation := strings.TrimRight(string(entry.Operation[:]), "\x00")
		path := strings.TrimRight(string(entry.Path[:]), "\x00")
		content := strings.TrimRight(string(entry.Content[:]), "\x00")
		date := strings.TrimRight(string(entry.Date[:]), "\x00")

		// Filtrar entradas con datos basura (todos caracteres nulos o patrones raros)
		if operation == "" || strings.Trim(operation, "\x001") == "" {
			continue
		}

		// Truncar strings largos para que quepan en la tabla
		if len(path) > 40 {
			path = path[:37] + "..."
		}
		if len(content) > 30 {
			content = content[:27] + "..."
		}

		fmt.Printf("%-15s %-40s %-30s %-20s\n", operation, path, content, date)
	}

	fmt.Println("═══════════════════════════════════════════════════════════════════")
}
