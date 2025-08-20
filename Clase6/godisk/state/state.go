package state

// MountedPartition representa una partición que ha sido cargada en memoria.
type MountedPartition struct {
	ID      string
	Path    string // Ruta al archivo de disco
	Name    string // Nombre de la partición
	Status  byte   // '1' para estado montada, '0' para no montada
	Correl  int    // Número correlativo asignado al montar
	Letter  rune   // Letra del disco
	PartNum int    // Número de partición en ese disco
	Size    int64  // Tamaño de la partición en bytes (CAMPO AÑADIDO)
	Start   int64  // Byte de inicio de la partición en el disco (CAMPO AÑADIDO)
}

// GlobalMountedPartitions es la lista en memoria de todas las particiones montadas.
var GlobalMountedPartitions []MountedPartition

func GetMountedPartitions() []MountedPartition {
	// Devuelve directamente la lista global.
	return GlobalMountedPartitions
}

// GetMountedPartitionByID busca en la lista global una partición por su ID.
// Devuelve la partición encontrada y un booleano 'true' si la encontró.
// Si no la encuentra, devuelve una estructura vacía y 'false'.
func GetMountedPartitionByID(id string) (MountedPartition, bool) {
	// Recorre la lista de particiones montadas.
	for _, p := range GlobalMountedPartitions {
		// Si el ID coincide, devuelve la partición y 'true'.
		if p.ID == id {
			return p, true
		}
	}
	// Si el bucle termina sin encontrar nada, devuelve valores vacíos y 'false'.
	return MountedPartition{}, false
}
