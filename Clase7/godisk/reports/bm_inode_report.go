package reports

import (
	"fmt"
	"godisk/structs"
	"godisk/utils"
	"os"
	"strings"
)

// GenerateBmInodeReport crea el contenido de texto para el reporte del bitmap de inodos.
func GenerateBmInodeReport(file *os.File, sb structs.Superblock) string {
	var content strings.Builder

	// 1. Leer el bitmap de inodos completo desde el disco.
	inodeBitmap, err := utils.ReadBytes(file, int64(sb.S_bm_inode_start), int64(sb.S_inodes_count))
	if err != nil {
		return fmt.Sprintf("Error al leer el bitmap de inodos: %v", err)
	}

	// 2. Recorrer el bitmap y formatear la salida.
	for i, status := range inodeBitmap {
		// Escribe el carácter directamente. Como ahora el disco tiene '0' y '1', esto funcionará.
		content.WriteByte(byte(status))
		content.WriteString(" ")

		// 3. Si hemos escrito 20 registros, añadir un salto de línea.
		if (i+1)%20 == 0 && i+1 < len(inodeBitmap) {
			content.WriteString("\n")
		}
	}

	return content.String()
}
