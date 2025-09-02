package reports

import (
	"encoding/binary"
	"fmt"
	"godisk/structs"
	"godisk/utils"
	"os"
	"sort"
	"strings"
)

// diskSegment es una estructura auxiliar para representar un trozo del disco.
type diskSegment struct {
	Start int64
	Size  int64
	Type  string // tipo de estructura: MBR, Primaria, Extendida, EBR, Lógica, Libre
	Name  string
}

// getExtendedSubLabel genera una tabla anidada para el contenido de una partición extendida.
func getExtendedSubLabel(file *os.File, extended structs.Partition, totalDiskSize int64) string {
	var subSegments []diskSegment
	currentPos := extended.Part_start
	ebrSize := int64(binary.Size(structs.EBR{}))

	//  Recorre la cadena de EBRs para recolectar EBRs y particiones lógicas.
	for {
		ebr, err := utils.ReadEBR(file, currentPos)
		if err != nil || ebr.Part_status == '0' {
			break // Si no se puede leer o está inactivo, se detiene.
		}

		// Añade el EBR como un segmento.
		subSegments = append(subSegments, diskSegment{Start: currentPos, Size: ebrSize, Type: "EBR"})
		// Añade la partición lógica como otro segmento.
		subSegments = append(subSegments, diskSegment{
			Start: ebr.Part_start,
			Size:  ebr.Part_s,
			Type:  "Lógica",
			Name:  strings.TrimRight(string(ebr.Part_name[:]), "\x00"),
		})

		if ebr.Part_next == -1 {
			break // Fin de la cadena.
		}
		currentPos = ebr.Part_next
	}

	// 2. Ordena los segmentos internos y calcular los espacios libres.
	sort.Slice(subSegments, func(i, j int) bool { return subSegments[i].Start < subSegments[j].Start })

	var finalSubSegments []diskSegment
	lastEnd := extended.Part_start
	for _, seg := range subSegments {
		if seg.Start > lastEnd {
			finalSubSegments = append(finalSubSegments, diskSegment{Type: "Libre", Size: seg.Start - lastEnd})
		}
		finalSubSegments = append(finalSubSegments, seg)
		lastEnd = seg.Start + seg.Size
	}
	if lastEnd < extended.Part_start+extended.Part_s {
		finalSubSegments = append(finalSubSegments, diskSegment{Type: "Libre", Size: (extended.Part_start + extended.Part_s) - lastEnd})
	}

	// 3. Construye la tabla HTML-like para la partición extendida.
	var subLabel strings.Builder
	subLabel.WriteString("<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"2\"><TR><TD>Extendida</TD></TR><TR>")
	for _, seg := range finalSubSegments {
		percentage := float64(seg.Size) * 100.0 / float64(totalDiskSize)
		label := fmt.Sprintf("%s<BR/>%.2f%% del Disco", seg.Type, percentage)
		subLabel.WriteString(fmt.Sprintf("<TD BGCOLOR=\"#FADBD8\">%s</TD>", label))
	}
	subLabel.WriteString("</TR></TABLE>")
	return subLabel.String()
}

// GenerateDiskReport crea el contenido DOT para el reporte gráfico del disco.
func GenerateDiskReport(file *os.File, mbr *structs.MBR) string {
	var segments []diskSegment
	mbrSize := int64(utils.Sizeof(mbr))

	// 1. Añadir el MBR y las particiones primarias/extendida.
	segments = append(segments, diskSegment{Start: 0, Size: mbrSize, Type: "MBR"})
	var extendedPartition structs.Partition
	hasExtended := false
	for _, p := range mbr.Mbr_partitions {
		if p.Part_status == '1' {
			segType := "Primaria"
			if p.Part_type == 'e' || p.Part_type == 'E' {
				segType = "Extendida"
				extendedPartition = p
				hasExtended = true
			}
			segments = append(segments, diskSegment{
				Start: p.Part_start,
				Size:  p.Part_s,
				Type:  segType,
			})
		}
	}

	// 2. Ordenar y calcular espacios libres.
	sort.Slice(segments, func(i, j int) bool { return segments[i].Start < segments[j].Start })
	var finalSegments []diskSegment
	lastEnd := int64(0)
	for _, seg := range segments {
		if seg.Start > lastEnd {
			finalSegments = append(finalSegments, diskSegment{Type: "Libre", Size: seg.Start - lastEnd})
		}
		finalSegments = append(finalSegments, seg)
		lastEnd = seg.Start + seg.Size
	}
	if lastEnd < mbr.Mbr_tamano {
		finalSegments = append(finalSegments, diskSegment{Type: "Libre", Size: mbr.Mbr_tamano - lastEnd})
	}

	// 3. Construir el string DOT final.
	var dot strings.Builder
	dot.WriteString("digraph Disk_Report {\n")
	dot.WriteString("    node [shape=plaintext];\n")
	dot.WriteString("    disk [label=<\n")
	dot.WriteString("        <TABLE BORDER=\"1\" CELLBORDER=\"0\" CELLSPACING=\"0\"><TR>\n")

	for _, seg := range finalSegments {
		// Si el segmento es la partición extendida, genera su tabla anidada.
		if seg.Type == "Extendida" && hasExtended {
			subLabel := getExtendedSubLabel(file, extendedPartition, mbr.Mbr_tamano)
			dot.WriteString(fmt.Sprintf("                <TD BORDER=\"1\">%s</TD>\n", subLabel))
		} else { // Para MBR, primarias y libres, crea una celda simple.
			percentage := float64(seg.Size) * 100.0 / float64(mbr.Mbr_tamano)
			label := fmt.Sprintf("%s<BR/>%.2f%% del Disco", seg.Type, percentage)
			if seg.Type == "MBR" {
				label = "MBR" // El MBR no muestra porcentaje.
			}
			dot.WriteString(fmt.Sprintf("                <TD BORDER=\"1\" BGCOLOR=\"#AED6F1\">%s</TD>\n", label))
		}
	}

	dot.WriteString("            </TR></TABLE>>];\n")
	dot.WriteString("}\n")
	return dot.String()
}
