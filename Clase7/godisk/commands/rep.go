package commands

import (
	"fmt"
	"godisk/reports"
	"godisk/state"
	"godisk/utils"
	"os"
	"path/filepath"
	"strings"
)

// ExecuteRep contiene la lógica para generar los diferentes reportes.
func ExecuteRep(name, path, id string) {
	mountedPartition, found := state.GetMountedPartitionByID(id)
	if !found {
		fmt.Printf("Error: no se encontró una partición montada con el ID '%s'.\n", id)
		return
	}

	file, err := os.Open(mountedPartition.Path)
	if err != nil {
		fmt.Printf("Error al abrir el disco '%s': %v\n", mountedPartition.Path, err)
		return
	}
	defer file.Close()

	var reportContent string
	isTextReport := false

	// Usar un switch para determinar qué reporte generar.
	switch strings.ToLower(name) {

	case "disk":
		mbr, err := utils.ReadMBR(file)
		if err != nil {
			fmt.Printf("Error al leer MBR: %v\n", err)
			return
		}
		// CORREGIDO: Pasar un puntero a mbr en lugar de un valor
		reportContent = reports.GenerateDiskReport(file, &mbr)

	case "bm_inode":
		sb, err := utils.ReadSuperblock(file, mountedPartition.Start)
		if err != nil {
			fmt.Printf("Error al leer el superbloque: %v\n", err)
			return
		}
		reportContent = reports.GenerateBmInodeReport(file, sb)
		isTextReport = true // Marcar como reporte de texto

	default:
		fmt.Printf("Error: el reporte con el nombre '%s' no es reconocido.\n", name)
		return
	}

	// Crear el directorio si no existe
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error al crear el directorio del reporte: %v\n", err)
		return
	}

	// Generar el archivo final
	if isTextReport {
		// Para reportes de texto, simplemente se escribe el contenido.
		err = os.WriteFile(path, []byte(reportContent), 0644)
		if err != nil {
			fmt.Printf("Error al escribir el reporte de texto: %v\n", err)
		} else {
			fmt.Printf("Reporte '%s' generado exitosamente en '%s'.\n", name, path)
		}
	} else {
		// Para reportes gráficos, se usa Graphviz.
		err = utils.GenerateReport(reportContent, path)
		if err != nil {
			fmt.Printf("Error al generar el reporte gráfico: %v\n", err)
		} else {
			fmt.Printf("Reporte '%s' generado exitosamente en '%s'.\n", name, path)
		}
	}
}
