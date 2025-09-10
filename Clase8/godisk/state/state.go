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
	Size    int64  // Tamaño de la partición en bytes
	Start   int64  // Byte de inicio de la partición en el disco
}

// User representa la información de la sesión de un usuario logueado.
type User struct {
	IsLoggedIn bool
	Id         int32  // ID del usuario
	Group      int32  // ID del grupo principal del usuario
	Username   string // Nombre de usuario
	MountedId  string // ID de la partición en la que está trabajando
}

// CurrentUser es la variable global que mantiene la sesión activa.
var CurrentUser User

// GlobalMountedPartitions es la lista en memoria de todas las particiones montadas.
var GlobalMountedPartitions []MountedPartition

// GetMountedPartitions devuelve la lista de particiones montadas.
func GetMountedPartitions() []MountedPartition {
	return GlobalMountedPartitions
}

// GetMountedPartitionByID busca una partición por su ID.
func GetMountedPartitionByID(id string) (MountedPartition, bool) {
	for _, p := range GlobalMountedPartitions {
		if p.ID == id {
			return p, true
		}
	}
	return MountedPartition{}, false
}
