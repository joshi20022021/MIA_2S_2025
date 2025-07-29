package commands

import (
	"encoding/binary"
	"fmt"
	"godisk/structs"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// esta funcion contiene la lógica para crear el disco.
func ExecuteMkdisk(size int, unit string, fit string, path string) {
	//Validar y calcular el tamaño real en bytes
	var diskSize int64

	// Validación completa del parámetro unit
	unit = strings.ToUpper(unit)
	if unit == "K" {
		diskSize = int64(size) * 1024
	} else if unit == "M" || unit == "" {
		diskSize = int64(size) * 1024 * 1024
	} else {
		fmt.Printf("Error: valor '%s' no válido para -unit. Use K o M.\n", unit)
		return
	}

	if diskSize <= 0 {
		fmt.Println("Error: el parámetro -size debe ser mayor a cero.")
		return
	}

	// Validación del parámetro fit
	var fitByte byte
	fit = strings.ToUpper(fit)
	if fit == "BF" {
		fitByte = 'b'
	} else if fit == "WF" {
		fitByte = 'w'
	} else if fit == "FF" || fit == "" {
		fitByte = 'f'
	} else {
		fmt.Printf("Error: valor '%s' no válido para -fit. Use BF, FF o WF.\n", fit)
		return
	}

	// Verificar la extensión .mia
	if !strings.HasSuffix(strings.ToLower(path), ".mia") {
		path += ".mia" // Añadir la extensión si no está
	}

	// 2. Crear las carpetas si no existen
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error al crear directorios: %v\n", err)
		return
	}

	// 3. Crear el archivo binario
	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error al crear el archivo: %v\n", err)
		return
	}
	defer file.Close()

	// 4. Llenar el archivo con ceros binarios
	chunk := make([]byte, 1024)
	for i := int64(0); i < diskSize/1024; i++ {
		if _, err := file.Write(chunk); err != nil {
			fmt.Printf("Error al escribir en el archivo: %v\n", err)
			return
		}
	}
	if err := file.Truncate(diskSize); err != nil {
		fmt.Printf("Error al truncar el archivo: %v\n", err)
		return
	}

	// 5. Crear el MBR
	rand.Seed(time.Now().UnixNano())
	diskSignature := rand.Int63()

	mbr := structs.NewMBR(diskSize, fitByte, diskSignature)

	// 6. Escribir el MBR al inicio del archivo
	file.Seek(0, 0) // Moverse al byte 0
	if err := binary.Write(file, binary.LittleEndian, &mbr); err != nil {
		fmt.Printf("Error al escribir el MBR: %v\n", err)
		return
	}

	fmt.Printf("Disco creado exitosamente en: %s\n", path)
	fmt.Printf("Tamaño: %d bytes, Firma: %d\n", mbr.Mbr_tamano, mbr.Mbr_dsk_signature)
}
