package commands

import (
	"encoding/binary"
	"fmt"
	"godisk/state"
	"godisk/structs"
	"math"
	"os"
	"strings"
	"time"
	"unsafe"
)

// ExecuteMkfs formatea una partición con un sistema de archivos.
// Recibe el ID de la partición montada y los tipos de formato y sistema de archivos.
func ExecuteMkfs(id, formatType, fsType string) {
	// VALIDACIÓN DE PARÁMETROS ---
	// aquí se asegura que el tipo de formato sea 'full'.
	if strings.ToLower(formatType) != "full" {
		fmt.Printf("Advertencia: tipo de formateo '%s' no reconocido. Se usará 'full'.\n", formatType)
		formatType = "full"
	}

	// --- BÚSQUEDA DE LA PARTICIÓN MONTADA ---
	// Llama a la función del paquete 'state' para obtener la información de la partición
	// que corresponde al ID proporcionado por el usuario.
	mountedPartition, found := state.GetMountedPartitionByID(id)
	if !found {
		// Si 'found' es false, la partición no está montada y no se puede continuar.
		fmt.Printf("Error: No se encontró una partición montada con el id '%s'.\n", id)
		return
	}

	fmt.Printf("Iniciando formateo para la partición %s en %s.\n", mountedPartition.Name, mountedPartition.Path)

	sizeOfSuperblock := int64(binary.Size(structs.Superblock{}))
	sizeOfInode := int64(binary.Size(structs.Inode{}))
	sizeOfBlock := int64(binary.Size(structs.FileBlock{}))

	// --- CÁLCULO DEL NÚMERO DE INODOS ---
	// Se calcula el número de inodos 'n' que caben en la partición.
	availableSpace := float64(mountedPartition.Size - sizeOfSuperblock)
	structureUnitSize := float64(sizeOfInode + (3 * sizeOfBlock))

	// --- INICIO DE BLOQUE DE DEPURACIÓN ---
	fmt.Println("------------------- particion INFO -------------------")
	fmt.Printf("Tamaño de la Partición (mountedPartition.Size): %d bytes\n", mountedPartition.Size)
	fmt.Printf("Tamaño del Superbloque (sizeOfSuperblock):      %d bytes\n", sizeOfSuperblock)
	fmt.Printf("Tamaño del Inodo (sizeOfInode):                 %d bytes\n", sizeOfInode)
	fmt.Printf("Tamaño del Bloque (sizeOfBlock):                %d bytes\n", sizeOfBlock)
	fmt.Println("--------------------------------------------------")
	fmt.Printf("Espacio Disponible (availableSpace):              %.f bytes\n", availableSpace)
	fmt.Printf("Tamaño de Unidad (structureUnitSize):             %.f bytes\n", structureUnitSize)
	fmt.Println("--------------------------------------------------")
	// --- FIN DE BLOQUE ---

	if structureUnitSize <= 0 {
		fmt.Println("Error: El tamaño de las estructuras del sistema de archivos es cero o negativo.")
		return
	}

	n := math.Floor(availableSpace / structureUnitSize)
	fmt.Printf("Número de Inodos Calculado (n):                 %.f\n", n) // Imprimir n también
	fmt.Println("--------------------------------------------------")

	// --- VALIDACIÓN DE ESPACIO ---
	// Si n es menor o igual a 0, no hay espacio suficiente en la partición para crear el sistema de archivos.
	if n <= 2 { // Se necesitan al menos 3 inodos (raíz, users.txt, y uno libre)
		fmt.Println("Error: Espacio insuficiente en la partición para crear el sistema de archivos.")
		return
	}

	// --- 4. CREACIÓN DEL SUPERBLOQUE EN MEMORIA ---
	// Se crea una instancia del Superbloque y se llena con la información calculada.
	var superbloque structs.Superblock
	superbloque.S_filesystem_type = 2 // 2 para ext2
	if fsType == "3fs" {
		superbloque.S_filesystem_type = 3
	}
	superbloque.S_inodes_count = int32(n)
	superbloque.S_blocks_count = 3 * int32(n)
	superbloque.S_free_blocks_count = 3 * int32(n)
	superbloque.S_free_inodes_count = int32(n)
	superbloque.S_mtime = time.Now().Unix()
	superbloque.S_umtime = time.Now().Unix()
	superbloque.S_mnt_count = 1
	superbloque.S_magic = 0xEF53
	superbloque.S_inode_size = int32(sizeOfInode)
	superbloque.S_block_size = int32(sizeOfBlock)

	// --- 5. CÁLCULO DE PUNTEROS DE INICIO ---
	// Se calculan las posiciones exactas (offsets en bytes) donde comenzará cada sección.
	partitionStart := mountedPartition.Start
	superbloque.S_bm_inode_start = int32(partitionStart + sizeOfSuperblock)
	superbloque.S_bm_block_start = superbloque.S_bm_inode_start + superbloque.S_inodes_count
	superbloque.S_inode_start = superbloque.S_bm_block_start + superbloque.S_blocks_count
	superbloque.S_block_start = superbloque.S_inode_start + (superbloque.S_inodes_count * superbloque.S_inode_size)
	// --- 6. APERTURA DEL ARCHIVO DE DISCO ---
	// Se abre el archivo del disco en modo lectura/escritura para poder modificarlo.
	file, err := os.OpenFile(mountedPartition.Path, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error al abrir el disco:", err)
		return
	}
	defer file.Close() // 'defer' asegura que el archivo se cierre al final de la función.

	// --- 7. ESCRITURA DEL SUPERBLOQUE EN DISCO ---
	// Se mueve el puntero del archivo al inicio de la partición.
	file.Seek(partitionStart, 0)
	// binary.Write convierte la struct 'superbloque' a su representación en bytes y la escribe en el archivo.
	if err := binary.Write(file, binary.BigEndian, &superbloque); err != nil {
		fmt.Println("Error al escribir el superbloque:", err)
		return
	}
	fmt.Println("Superbloque creado y escrito.")

	// --- 8. ESCRITURA DE BITMAPS Y BLOQUES (FORMATEO FULL) ---
	// Se crean slices de bytes (arrays) para los bitmaps, inicializados en cero.
	bmInode := make([]byte, superbloque.S_inodes_count)
	for i := range bmInode {
		bmInode[i] = '0' // CORRECCIÓN: Inicializar con el carácter '0'
	}
	bmBlock := make([]byte, superbloque.S_blocks_count)
	for i := range bmBlock {
		bmBlock[i] = '0' // CORRECCIÓN: Inicializar con el carácter '0'
	}

	if strings.ToLower(formatType) == "full" {
		fmt.Println("Realizando formateo completo (full)...")
		// Se escribe el bitmap de inodos (lleno de '0') en su posición.
		file.Seek(int64(superbloque.S_bm_inode_start), 0)
		binary.Write(file, binary.BigEndian, &bmInode)

		// Se escribe el bitmap de bloques (lleno de '0') en su posición.
		file.Seek(int64(superbloque.S_bm_block_start), 0)
		binary.Write(file, binary.BigEndian, &bmBlock)

		// Se llenan las tablas de inodos y bloques con ceros para una limpieza completa.
		emptyInode := structs.Inode{}
		file.Seek(int64(superbloque.S_inode_start), 0)
		for i := 0; i < int(superbloque.S_inodes_count); i++ {
			binary.Write(file, binary.BigEndian, &emptyInode)
		}

		emptyBlock := make([]byte, sizeOfBlock)
		file.Seek(int64(superbloque.S_block_start), 0)
		for i := 0; i < int(superbloque.S_blocks_count); i++ {
			binary.Write(file, binary.BigEndian, &emptyBlock)
		}
	}
	fmt.Println("Bitmaps y bloques inicializados.")

	// --- 9. CREACIÓN DEL SISTEMA DE ARCHIVOS RAÍZ Y USERS.TXT ---
	// Se crea el inodo para el directorio raíz ("/") en memoria (Inodo 0).
	var rootInode structs.Inode
	rootInode.I_uid = 1                                            // Propietario: root
	rootInode.I_gid = 1                                            // Grupo: root
	rootInode.I_size = int32(unsafe.Sizeof(structs.FolderBlock{})) // Ocupa un bloque de carpeta
	rootInode.I_atime = time.Now().Unix()
	rootInode.I_ctime = time.Now().Unix()
	rootInode.I_mtime = time.Now().Unix()
	rootInode.I_type = 0   // 0 para carpeta
	rootInode.I_perm = 664 // Permisos de lectura/escritura para propietario/grupo, lectura para otros
	for i := range rootInode.I_block {
		rootInode.I_block[i] = -1
	}
	rootInode.I_block[0] = 0 // El primer puntero directo apunta al bloque de datos 0.

	// Se crea el inodo para el archivo "users.txt" (Inodo 1).
	usersContent := "1,G,root\n1,U,root,root,123\n"
	var usersInode structs.Inode
	usersInode.I_uid = 1
	usersInode.I_gid = 1
	usersInode.I_size = int32(len(usersContent)) // Tamaño exacto del contenido
	usersInode.I_atime = time.Now().Unix()
	usersInode.I_ctime = time.Now().Unix()
	usersInode.I_mtime = time.Now().Unix()
	usersInode.I_type = 1 // 1 para archivo
	usersInode.I_perm = 664
	for i := range usersInode.I_block {
		usersInode.I_block[i] = -1
	}
	usersInode.I_block[0] = 1 // El primer puntero directo apunta al bloque de datos 1.

	// Se crea el bloque de carpeta para la raíz (Bloque 0).
	var rootFolderBlock structs.FolderBlock
	// Inicializar todas las entradas a -1 para indicar que están vacías
	for i := range rootFolderBlock.B_content {
		rootFolderBlock.B_content[i].B_inodo = -1
	}
	// Entrada para sí mismo "."
	copy(rootFolderBlock.B_content[0].B_name[:], ".")
	rootFolderBlock.B_content[0].B_inodo = 0
	// Entrada para el padre ".."
	copy(rootFolderBlock.B_content[1].B_name[:], "..")
	rootFolderBlock.B_content[1].B_inodo = 0 // La raíz es su propio padre
	// Entrada para "users.txt"
	copy(rootFolderBlock.B_content[2].B_name[:], "users.txt")
	rootFolderBlock.B_content[2].B_inodo = 1 // Apunta al inodo 1

	// Se crea el bloque de contenido para "users.txt" (Bloque 1).
	var usersFileBlock structs.FileBlock
	copy(usersFileBlock.B_content[:], usersContent)

	// --- 10. ESCRITURA DE ESTRUCTURAS INICIALES ---
	// Escribir inodo raíz (inodo 0)
	file.Seek(int64(superbloque.S_inode_start), 0)
	binary.Write(file, binary.BigEndian, &rootInode)
	// Escribir inodo de users.txt (inodo 1)
	file.Seek(int64(superbloque.S_inode_start)+int64(superbloque.S_inode_size), 0)
	binary.Write(file, binary.BigEndian, &usersInode)

	// Escribir bloque de carpeta raíz (bloque 0)
	file.Seek(int64(superbloque.S_block_start), 0)
	binary.Write(file, binary.BigEndian, &rootFolderBlock)
	// Escribir bloque de contenido de users.txt (bloque 1)
	file.Seek(int64(superbloque.S_block_start)+int64(superbloque.S_block_size), 0)
	binary.Write(file, binary.BigEndian, &usersFileBlock)

	// --- 11. ACTUALIZACIÓN DE ESTRUCTURAS DE CONTROL ---
	// Marcar inodos 0 y 1 como usados ('1') en el bitmap.
	file.Seek(int64(superbloque.S_bm_inode_start), 0)
	// CORRECCIÓN: Escribir los caracteres '1', '1'
	binary.Write(file, binary.BigEndian, []byte{'1', '1'})

	// Marcar bloques 0 y 1 como usados ('1') en el bitmap.
	file.Seek(int64(superbloque.S_bm_block_start), 0)
	// CORRECCIÓN: Escribir los caracteres '1', '1'
	binary.Write(file, binary.BigEndian, []byte{'1', '1'})

	// Se actualizan los contadores y punteros a libres en la copia en memoria del superbloque.
	superbloque.S_free_inodes_count -= 2
	superbloque.S_free_blocks_count -= 2
	superbloque.S_first_ino = 2 // El siguiente inodo libre es ahora el 2.
	superbloque.S_first_blo = 2 // El siguiente bloque libre es ahora el 2.

	// Finalmente, se reescribe el superbloque completo en el disco con los valores actualizados.
	file.Seek(partitionStart, 0)
	binary.Write(file, binary.BigEndian, &superbloque)

	fmt.Println("Sistema de archivos creado exitosamente en la partición, incluyendo users.txt.")
}
