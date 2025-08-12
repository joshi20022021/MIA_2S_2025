package structs

// EBR representa un Extended Boot Record.
// Se usa para gestionar particiones lógicas dentro de una partición extendida.
type EBR struct {
	// Part_status: Indica si la partición está activa ('1') o no ('0')
	Part_status byte
	// Part_fit: Tipo de ajuste para la partición lógica ('B', 'F', 'W')
	Part_fit byte
	// Part_start: Byte donde inicia esta partición lógica
	Part_start int64
	// Part_s: Tamaño en bytes de esta partición lógica
	Part_s int64
	// Part_next: Apuntador al byte donde se encuentra el próximo EBR.
	// Es -1 si no hay más particiones lógicas.
	Part_next int64
	// Part_name: Nombre de la partición lógica
	Part_name [16]byte
}
