package reports

import (
	"fmt"
	"godisk/structs"
	"godisk/utils"
	"os"
	"strings"
	"time"
)

// generateLogicalPartitionRows es una función auxiliar recursiva para añadir filas de EBRs.
func generateLogicalPartitionRows(file *os.File, start int64, dot *strings.Builder) {
	if start == -1 {
		return
	}

	ebr, err := utils.ReadEBR(file, start)
	if err != nil {
		return
	}

	// Solo procesa la partición lógica si está activa.
	if ebr.Part_status == '1' {
		cleanEbrName := strings.TrimRight(string(ebr.Part_name[:]), "\x00")

		// Añade las filas para esta partición lógica.
		dot.WriteString("            <TR><TD COLSPAN=\"2\" BGCOLOR=\"#D2B4DE\"><B>EBR (Partición Lógica)</B></TD></TR>\n")
		dot.WriteString(fmt.Sprintf("            <TR><TD>part_status</TD><TD>%c</TD></TR>\n", ebr.Part_status))
		dot.WriteString(fmt.Sprintf("            <TR><TD>part_fit</TD><TD>%c</TD></TR>\n", ebr.Part_fit))
		dot.WriteString(fmt.Sprintf("            <TR><TD>part_start</TD><TD>%d</TD></TR>\n", ebr.Part_start))
		dot.WriteString(fmt.Sprintf("            <TR><TD>part_s</TD><TD>%d</TD></TR>\n", ebr.Part_s))
		dot.WriteString(fmt.Sprintf("            <TR><TD>part_next</TD><TD>%d</TD></TR>\n", ebr.Part_next))
		dot.WriteString(fmt.Sprintf("            <TR><TD>part_name</TD><TD>%s</TD></TR>\n", cleanEbrName))
	}

	// Llamada recursiva para el siguiente EBR en la cadena.
	generateLogicalPartitionRows(file, ebr.Part_next, dot)
}

// GenerateMbrReport crea el contenido DOT para el reporte de MBR, incluyendo EBRs.
func GenerateMbrReport(file *os.File, mbr *structs.MBR) string {
	var dot strings.Builder
	dot.WriteString("digraph MBR_Report {\n")
	dot.WriteString("    node [shape=plaintext];\n\n")

	dot.WriteString("    mbr_table [label=<\n")
	dot.WriteString("        <TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\">\n")
	// Datos del MBR.
	dot.WriteString("            <TR><TD COLSPAN=\"2\" BGCOLOR=\"#AED6F1\"><B>MBR</B></TD></TR>\n")
	dot.WriteString(fmt.Sprintf("            <TR><TD>mbr_tamano</TD><TD>%d</TD></TR>\n", mbr.Mbr_tamano))
	dot.WriteString(fmt.Sprintf("            <TR><TD>mbr_fecha_creacion</TD><TD>%s</TD></TR>\n", time.Unix(mbr.Mbr_fecha_creacion, 0).Format("2006-01-02 15:04:05")))
	dot.WriteString(fmt.Sprintf("            <TR><TD>mbr_dsk_signature</TD><TD>%d</TD></TR>\n", mbr.Mbr_dsk_signature))

	// Itera sobre las 4 particiones del MBR.
	for i, p := range mbr.Mbr_partitions {
		if p.Part_status == '1' {
			cleanPartName := strings.TrimRight(string(p.Part_name[:]), "\x00")
			partHeader := fmt.Sprintf("Partición %d", i+1)
			if p.Part_type == 'e' || p.Part_type == 'E' {
				partHeader += " (Extendida)"
			} else {
				partHeader += " (Primaria)"
			}

			dot.WriteString(fmt.Sprintf("            <TR><TD COLSPAN=\"2\" BGCOLOR=\"#A9DFBF\"><B>%s</B></TD></TR>\n", partHeader))
			dot.WriteString(fmt.Sprintf("            <TR><TD>part_status</TD><TD>%c</TD></TR>\n", p.Part_status))
			dot.WriteString(fmt.Sprintf("            <TR><TD>part_type</TD><TD>%c</TD></TR>\n", p.Part_type))
			dot.WriteString(fmt.Sprintf("            <TR><TD>part_fit</TD><TD>%c</TD></TR>\n", p.Part_fit))
			dot.WriteString(fmt.Sprintf("            <TR><TD>part_start</TD><TD>%d</TD></TR>\n", p.Part_start))
			dot.WriteString(fmt.Sprintf("            <TR><TD>part_s</TD><TD>%d</TD></TR>\n", p.Part_s))
			dot.WriteString(fmt.Sprintf("            <TR><TD>part_name</TD><TD>%s</TD></TR>\n", cleanPartName))

			// Si la partición es Extendida, llama a la función para insertar los EBRs.
			if p.Part_type == 'e' || p.Part_type == 'E' {
				generateLogicalPartitionRows(file, p.Part_start, &dot)
			}
		}
	}

	dot.WriteString("        </TABLE>>];\n")
	dot.WriteString("}\n")
	return dot.String()
}
