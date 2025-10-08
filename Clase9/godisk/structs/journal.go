package structs

// JournalEntry representa una entrada del journal del sistema de archivos
type JournalEntry struct {
	Operation [10]byte  // Tipo de operación (mkdir, mkfile, etc.)
	Path      [100]byte // Ruta del archivo/directorio
	Content   [100]byte // Contenido o descripción
	Date      [20]byte  // Fecha y hora de la operación
}
