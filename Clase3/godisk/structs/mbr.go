package structs

import "time"

// Partition representa la estructura de una partición dentro del MBR.
type Partition struct {
	Part_status byte
	Part_type   byte
	Part_fit    byte
	Part_start  int64
	Part_s      int64
	Part_name   [16]byte
}

// MBR representa el Master Boot Record del disco.
type MBR struct {
	Mbr_tamano         int64
	Mbr_fecha_creacion int64
	Mbr_dsk_signature  int64
	Dsk_fit            byte
	Mbr_partitions     [4]Partition
}

// NewMBR es una función para crear un MBR con valores iniciales.
func NewMBR(size int64, fit byte, signature int64) MBR {
	var mbr MBR
	mbr.Mbr_tamano = size
	mbr.Mbr_fecha_creacion = time.Now().Unix()
	mbr.Mbr_dsk_signature = signature
	mbr.Dsk_fit = fit

	for i := 0; i < 4; i++ {
		mbr.Mbr_partitions[i].Part_status = '0'
		mbr.Mbr_partitions[i].Part_start = -1
	}

	return mbr
}
