package structs

// ContentEntry representa una entrada en un bloque de directorio (nombre -> inodo).
type ContentEntry struct {
	B_name  [12]byte // Nombre del archivo o carpeta
	B_inodo int32    // Apuntador al inodo
}

// FolderBlock es la estructura para un bloque de directorio.
type FolderBlock struct {
	B_content [4]ContentEntry
}

// FileBlock es la estructura para un bloque de contenido de archivo.
type FileBlock struct {
	B_content [64]byte
}
