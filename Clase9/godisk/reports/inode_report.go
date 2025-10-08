package reports

import (
	"fmt"
	"godisk/structs"
	"godisk/utils"
	"os"
	"strings"
	"time"
)

// GenerateInodeReport crea el contenido DOT para el reporte gráfico de inodos.
func GenerateInodeReport(file *os.File, sb structs.Superblock) string {
	var dot strings.Builder
	dot.WriteString("digraph Inode_Report {\n")
	dot.WriteString("  rankdir=TB;\n")
	dot.WriteString("  node [shape=plaintext];\n")
	dot.WriteString("  graph [pad=\"0.5\", nodesep=\"1\", ranksep=\"2\"];\n\n")

	// 1. Leer el bitmap de inodos completo
	inodeBitmap, err := utils.ReadBytes(file, int64(sb.S_bm_inode_start), int64(sb.S_inodes_count))
	if err != nil {
		return fmt.Sprintf("Error al leer el bitmap de inodos: %v", err)
	}

	foundInodes := false
	var usedInodes []structs.Inode
	var usedInodeIndices []int

	// --- PRIMERA PASADA: DIBUJAR TODOS LOS NODOS ---
	for i, status := range inodeBitmap {
		if status == '1' { // 1 significa que el inodo está en uso
			foundInodes = true

			inodeStart := int64(sb.S_inode_start) + int64(i)*int64(utils.Sizeof(structs.Inode{}))
			inode, err := utils.ReadInode(file, inodeStart)
			if err != nil {
				continue
			}

			// Guardar el inodo y su índice para la segunda pasada
			usedInodes = append(usedInodes, inode)
			usedInodeIndices = append(usedInodeIndices, i)

			// Generar la tabla DOT para este inodo
			dot.WriteString(generateInodeTable(i, inode))
		}
	}

	// --- segunda pasada: unir los nodos ---
	if foundInodes && len(usedInodeIndices) > 1 {
		// flecha entre el inodo raíz entre el índice 0 y el inodo del archivo `users.txt` índice 1
		dot.WriteString(fmt.Sprintf("  inode_%d -> inode_%d [label=\"users.txt\"];\n", usedInodeIndices[0], usedInodeIndices[1]))
	}

	// Si no se encuentra ningún inodo, añadir un nodo de texto para evitar un archivo vacío.
	if !foundInodes {
		dot.WriteString("  no_inodes [label=\"No se encontraron inodos en uso.\"];\n")
	}

	dot.WriteString("}\n")
	return dot.String()
}
func generateInodeTable(index int, inode structs.Inode) string {
	var table strings.Builder

	// Usar formato válido para Graphviz
	table.WriteString(fmt.Sprintf("  inode_%d [label=<\n", index))
	table.WriteString("    <TABLE BORDER=\"1\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"4\">\n")
	table.WriteString(fmt.Sprintf("      <TR><TD COLSPAN=\"2\" BGCOLOR=\"lightblue\"><B>Inodo %d</B></TD></TR>\n", index))

	// Información básica del inodo
	table.WriteString(fmt.Sprintf("      <TR><TD ALIGN=\"LEFT\">i_uid</TD><TD ALIGN=\"LEFT\">%d</TD></TR>\n", inode.I_uid))
	table.WriteString(fmt.Sprintf("      <TR><TD ALIGN=\"LEFT\">i_gid</TD><TD ALIGN=\"LEFT\">%d</TD></TR>\n", inode.I_gid))
	table.WriteString(fmt.Sprintf("      <TR><TD ALIGN=\"LEFT\">i_size</TD><TD ALIGN=\"LEFT\">%d</TD></TR>\n", inode.I_size))

	// Formatear las fechas de forma segura
	atime := time.Unix(inode.I_atime, 0).Format("2006-01-02 15:04:05")
	ctime := time.Unix(inode.I_ctime, 0).Format("2006-01-02 15:04:05")
	mtime := time.Unix(inode.I_mtime, 0).Format("2006-01-02 15:04:05")

	table.WriteString(fmt.Sprintf("      <TR><TD ALIGN=\"LEFT\">i_atime</TD><TD ALIGN=\"LEFT\">%s</TD></TR>\n", atime))
	table.WriteString(fmt.Sprintf("      <TR><TD ALIGN=\"LEFT\">i_ctime</TD><TD ALIGN=\"LEFT\">%s</TD></TR>\n", ctime))
	table.WriteString(fmt.Sprintf("      <TR><TD ALIGN=\"LEFT\">i_mtime</TD><TD ALIGN=\"LEFT\">%s</TD></TR>\n", mtime))

	// Mostrar los punteros a bloques unicamente solo los que están en uso
	for i, block := range inode.I_block {
		if block != -1 { // Solo mostrar punteros utilizados
			blockType := ""
			if i < 12 {
				blockType = "directo"
			} else if i == 12 {
				blockType = "indirecto"
			} else if i == 13 {
				blockType = "doble ind."
			} else if i == 14 {
				blockType = "triple ind."
			}
			table.WriteString(fmt.Sprintf("      <TR><TD ALIGN=\"LEFT\">i_block[%d] (%s)</TD><TD ALIGN=\"LEFT\">%d</TD></TR>\n", i, blockType, block))
		}
	}

	// Interpretar el tipo de inodo
	var inodeType string
	switch inode.I_type {
	case 1:
		inodeType = "Archivo"
	case 0:
		inodeType = "Directorio"
	default:
		inodeType = fmt.Sprintf("Desconocido (%d)", inode.I_type)
	}
	table.WriteString(fmt.Sprintf("      <TR><TD ALIGN=\"LEFT\">i_type</TD><TD ALIGN=\"LEFT\">%s</TD></TR>\n", inodeType))

	// Mostrar permisos en formato octal y descripción
	permsDesc := formatPermissions(int(inode.I_perm))
	table.WriteString(fmt.Sprintf("      <TR><TD ALIGN=\"LEFT\">i_perm</TD><TD ALIGN=\"LEFT\">%d (%s)</TD></TR>\n", inode.I_perm, permsDesc))

	table.WriteString("    </TABLE>\n")
	table.WriteString("  >];\n\n")

	return table.String()
}

// convierte permisos numéricos a descripción legible
func formatPermissions(perm int) string {
	if perm <= 0 {
		return "---"
	}

	// Convertir a octal y formatear
	owner := (perm / 100) % 10
	group := (perm / 10) % 10
	other := perm % 10

	permToStr := func(p int) string {
		str := ""
		if p&4 != 0 {
			str += "r"
		} else {
			str += "-"
		}
		if p&2 != 0 {
			str += "w"
		} else {
			str += "-"
		}
		if p&1 != 0 {
			str += "x"
		} else {
			str += "-"
		}
		return str
	}

	return permToStr(owner) + permToStr(group) + permToStr(other)
}
