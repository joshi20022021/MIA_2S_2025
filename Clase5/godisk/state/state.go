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
}

// GlobalMountedPartitions es la lista en memoria de todas las particiones montadas.
var GlobalMountedPartitions []MountedPartition
