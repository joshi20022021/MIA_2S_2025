package structs

import "time" // Paquete para manejar fechas y tiempo

// Partition representa la estructura de una partición dentro del MBR.
// Esta estructura define cómo se almacena la información de cada partición
type Partition struct {
	// Part_status: 1 byte que indica el estado de la partición
	// '0' = inactiva, '1' = activa
	Part_status byte

	// Part_type: 1 byte que indica el tipo de partición
	// 'p' = primaria, 'e' = extendida, 'l' = lógica
	Part_type byte

	// Part_fit: 1 byte que indica el tipo de ajuste para esta partición
	// 'b' = Best Fit, 'f' = First Fit, 'w' = Worst Fit
	Part_fit byte

	// Part_start: 8 bytes que indican la posición de inicio de la partición en el disco
	// Se mide en bytes desde el inicio del archivo
	Part_start int64

	// Part_s: 8 bytes que indican el tamaño de la partición
	// Se mide en bytes
	Part_s int64

	// Part_name: Array fijo de 16 bytes para el nombre de la partición
	// Se usa array fijo para tener un tamaño predecible en memoria
	Part_name [16]byte
}

// MBR representa el Master Boot Record del disco.
// Esta es la estructura principal que se escribe al inicio de cada disco virtual
type MBR struct {
	// Mbr_tamano: 8 bytes que almacenan el tamaño total del disco en bytes
	Mbr_tamano int64

	// Mbr_fecha_creacion: 8 bytes que almacenan el timestamp Unix de creación
	// Unix timestamp = segundos transcurridos desde el 1 de enero de 1970
	Mbr_fecha_creacion int64

	// Mbr_dsk_signature: 8 bytes con un número aleatorio que identifica únicamente este disco
	// Sirve como "firma" o "ID" único del disco
	Mbr_dsk_signature int64

	// Dsk_fit: 1 byte que indica el tipo de ajuste por defecto del disco
	// Este valor se usa cuando no se especifica ajuste para las particiones
	Dsk_fit byte

	// Mbr_partitions: Array fijo de exactamente 4 particiones
	// En los sistemas reales, el MBR solo puede contener 4 particiones primarias
	Mbr_partitions [4]Partition
}

// NewMBR es una función constructora para crear un MBR con valores iniciales.
// Recibe los parámetros necesarios y retorna un MBR completamente inicializado
func NewMBR(size int64, fit byte, signature int64) MBR {
	// Declara una variable de tipo MBR
	// En Go, las estructuras se inicializan automáticamente con valores cero
	var mbr MBR

	// Asigna el tamaño del disco al MBR
	mbr.Mbr_tamano = size

	// Obtiene el timestamp actual en formato Unix (segundos desde 1970)
	// time.Now() obtiene el tiempo actual, .Unix() lo convierte a timestamp
	mbr.Mbr_fecha_creacion = time.Now().Unix()

	// Asigna la firma única proporcionada
	mbr.Mbr_dsk_signature = signature

	// Asigna el tipo de ajuste por defecto del disco
	mbr.Dsk_fit = fit

	// Inicializa las 4 particiones con valores por defecto
	for i := 0; i < 4; i++ {
		// Marca cada partición como inactiva usando el carácter '0'
		mbr.Mbr_partitions[i].Part_status = '0'

		// Establece la posición de inicio como -1 para indicar que no está asignada
		// -1 es un valor especial que indica "sin asignar" o "no utilizada"
		mbr.Mbr_partitions[i].Part_start = -1

		// Los demás campos (Part_type, Part_fit, Part_s, Part_name) se inicializan
		// automáticamente con valores cero (0 para números, arrays vacíos para arrays)
	}

	// Retorna el MBR completamente inicializado
	// Se retorna por valor (copia) no por referencia
	return mbr
}
