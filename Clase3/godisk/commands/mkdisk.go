package commands

// Importaciones necesarias para la funcionalidad del comando mkdisk
import (
	"encoding/binary" // Para convertir estructuras Go a formato binario
	"fmt"             // Para imprimir mensajes en consola
	"godisk/structs"  // Nuestro paquete que contiene las definiciones de MBR y Partition
	"math/rand"       // Para generar números aleatorios (firma del disco)
	"os"              // Para operaciones del sistema operativo (crear archivos, directorios)
	"path/filepath"   // Para manipular rutas de archivos de forma segura entre plataformas
	"strings"         // Para manipular cadenas de texto
	"time"            // Para obtener tiempo actual y sembrar el generador de aleatorios
)

// ExecuteMkdisk contiene la lógica principal para crear un disco virtual.
// Esta función es exportada (empieza con mayúscula) para que pueda ser llamada desde otros paquetes
func ExecuteMkdisk(size int, unit string, fit string, path string) {

	// Declara variable para almacenar el tamaño final en bytes
	// Se usa int64 para soportar discos grandes (hasta 9 exabytes teóricamente)
	var diskSize int64

	// Validación completa del parámetro unit (unidad de medida)
	// Convierte a mayúsculas para hacer comparación insensible a mayúsculas/minúsculas
	unit = strings.ToUpper(unit)

	if unit == "K" {
		// Si la unidad es K (kilobytes), multiplica por 1024
		// 1 KB = 1024 bytes (sistema binario, no decimal)
		diskSize = int64(size) * 1024
	} else if unit == "M" || unit == "" {
		// Si la unidad es M (megabytes) o está vacía (valor por defecto)
		// 1 MB = 1024 * 1024 = 1,048,576 bytes
		diskSize = int64(size) * 1024 * 1024
	} else {
		// Si la unidad no es válida, mostrar error y terminar función
		fmt.Printf("Error: valor '%s' no válido para -unit. Use K o M.\n", unit)
		return // Termina la ejecución de la función
	}

	// Validación adicional: el tamaño debe ser positivo
	if diskSize <= 0 {
		fmt.Println("Error: el parámetro -size debe ser mayor a cero.")
		return
	}

	// Variable para almacenar el byte que representa el tipo de ajuste
	var fitByte byte

	// Convierte a mayúsculas para comparación insensible a mayúsculas/minúsculas
	fit = strings.ToUpper(fit)

	if fit == "BF" {
		// BF = Best Fit (mejor ajuste): busca el espacio libre más pequeño que pueda contener el archivo
		fitByte = 'b'
	} else if fit == "WF" {
		// WF = Worst Fit (peor ajuste): busca el espacio libre más grande disponible
		fitByte = 'w'
	} else if fit == "FF" || fit == "" {
		// FF = First Fit (primer ajuste): usa el primer espacio libre que encuentre
		// Este es el valor por defecto si no se especifica
		fitByte = 'f'
	} else {
		// Si el valor no es válido, mostrar error y terminar
		fmt.Printf("Error: valor '%s' no válido para -fit. Use BF, FF o WF.\n", fit)
		return
	}

	// Verifica si la ruta termina con la extensión .mia (insensible a mayúsculas)
	if !strings.HasSuffix(strings.ToLower(path), ".mia") {
		// Si no tiene la extensión, la añade automáticamente
		path += ".mia"
	}

	// filepath.Dir() extrae el directorio de la ruta completa
	// Por ejemplo: "/home/user/discos/disco1.mia" -> "/home/user/discos"
	dir := filepath.Dir(path)

	// os.MkdirAll() crea todos los directorios necesarios en la ruta
	// Es similar a "mkdir -p" en Unix: crea directorios padre si no existen
	// 0755 son los permisos: rwxr-xr-x (lectura/escritura/ejecución para owner, lectura/ejecución para otros)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error al crear directorios: %v\n", err)
		return
	}

	// os.Create() crea un nuevo archivo o trunca uno existente
	// Retorna un puntero al archivo y un error (patrón estándar de Go)
	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error al crear el archivo: %v\n", err)
		return
	}

	// defer asegura que el archivo se cierre cuando la función termine
	// Se ejecuta al final, sin importar cómo termine la función (éxito o error)
	defer file.Close()

	// Crea un slice de 1024 bytes (1 KB) inicializado con ceros
	// make() inicializa automáticamente con valores cero (0 para bytes)
	chunk := make([]byte, 1024)

	// Bucle para escribir chunks de 1 KB hasta llenar casi todo el disco
	// Se usa int64 para manejar discos grandes
	for i := int64(0); i < diskSize/1024; i++ {
		// file.Write() escribe el chunk al archivo
		// El "_" ignora el número de bytes escritos, solo nos interesa el error
		if _, err := file.Write(chunk); err != nil {
			fmt.Printf("Error al escribir en el archivo: %v\n", err)
			return
		}
	}

	// file.Truncate() ajusta el archivo al tamaño exacto deseado
	// Esto es necesario porque el bucle anterior podría dejar el archivo ligeramente más pequeño
	// (si diskSize no es múltiplo exacto de 1024)
	if err := file.Truncate(diskSize); err != nil {
		fmt.Printf("Error al truncar el archivo: %v\n", err)
		return
	}

	// Inicializa el generador de números aleatorios con el tiempo actual en nanosegundos
	// Esto asegura que cada ejecución produzca números diferentes
	rand.Seed(time.Now().UnixNano())

	// Genera un número aleatorio de 63 bits para usar como firma única del disco
	// rand.Int63() genera números positivos (no usa el bit de signo)
	diskSignature := rand.Int63()

	// Llama al constructor NewMBR para crear la estructura MBR
	// Pasa el tamaño del disco, el tipo de ajuste y la firma única
	mbr := structs.NewMBR(diskSize, fitByte, diskSignature)

	// file.Seek(0, 0) mueve el puntero de escritura al byte 0 (inicio del archivo)
	// Primer parámetro: offset (0 = inicio)
	// Segundo parámetro: whence (0 = desde el inicio del archivo)
	file.Seek(0, 0)

	// binary.Write() serializa la estructura MBR a formato binario y la escribe
	// - file: donde escribir
	// - binary.LittleEndian: orden de bytes (little endian es estándar en x86)
	// - &mbr: dirección de memoria de la estructura (se pasa por referencia)
	if err := binary.Write(file, binary.LittleEndian, &mbr); err != nil {
		fmt.Printf("Error al escribir el MBR: %v\n", err)
		return
	}

	// Informa al usuario que el disco se creó correctamente
	fmt.Printf("Disco creado exitosamente en: %s\n", path)

	// Muestra información relevante del disco creado
	// mbr.Mbr_tamano: tamaño total en bytes
	// mbr.Mbr_dsk_signature: firma única generada
	fmt.Printf("Tamaño: %d bytes, Firma: %d\n", mbr.Mbr_tamano, mbr.Mbr_dsk_signature)

	// La función termina aquí. El defer file.Close() se ejecuta automáticamente
}
