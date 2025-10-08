package fs

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"godisk/state"
	"godisk/structs"
	"godisk/utils"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FindInodeByPath navega el sistema de archivos para encontrar el inodo de una ruta específica.
func FindInodeByPath(file *os.File, sb structs.Superblock, path string) (structs.Inode, int32, error) {
	if path == "" {
		return structs.Inode{}, -1, errors.New("ruta vacía")
	}

	// Normalizar ruta: quitar espacios, comillas y limpiar
	cleanPath := strings.TrimSpace(path)
	cleanPath = strings.Trim(cleanPath, `"'`)
	cleanPath = filepath.Clean(cleanPath)

	if !strings.HasPrefix(cleanPath, "/") {
		return structs.Inode{}, -1, errors.New("la ruta debe ser absoluta (empezar con /)")
	}

	// caso raíz
	if cleanPath == "/" {
		root, err := ReadInode(file, sb, 0)
		return root, 0, err
	}

	parts := strings.Split(strings.Trim(cleanPath, "/"), "/")
	currentIndex := int32(0) // inodo raíz

	for _, part := range parts {
		if part == "" || part == "." {
			continue
		}

		currentInode, err := ReadInode(file, sb, currentIndex)
		if err != nil {
			return structs.Inode{}, -1, fmt.Errorf("error leyendo inodo %d: %v", currentIndex, err)
		}

		// si no es directorio, no podemos descender más (salvo si es el último componente y buscamos archivo)
		if currentInode.I_type != 0 {
			return structs.Inode{}, -1, fmt.Errorf("ruta inválida: '%s' no es un directorio intermedio", part)
		}

		// buscar la entrada en el bloque(s) del directorio actual
		nextIndex, err := findEntryInInode(file, sb, currentInode, part)
		if err != nil {
			return structs.Inode{}, -1, fmt.Errorf("no se encontró '%s' en la ruta '%s': %v", part, cleanPath, err)
		}
		currentIndex = nextIndex
	}

	finalInode, err := ReadInode(file, sb, currentIndex)
	if err != nil {
		return structs.Inode{}, -1, fmt.Errorf("error leyendo inodo final %d: %v", currentIndex, err)
	}
	return finalInode, currentIndex, nil
}

// ReadFileContent lee todos los bloques de datos de un inodo y devuelve su contenido.
func ReadFileContent(file *os.File, sb structs.Superblock, inode structs.Inode) ([]byte, error) {
	if inode.I_type != 1 { // 1 es para archivo
		return nil, errors.New("el inodo no corresponde a un archivo")
	}
	var content bytes.Buffer
	// solo leemos punteros directos
	for i := 0; i < 12 && inode.I_block[i] != -1; i++ {
		blockPtr := inode.I_block[i]
		fileBlock, err := ReadFileBlock(file, sb, blockPtr)
		if err != nil {
			return nil, err
		}
		// Asegurarse de no leer más allá del tamaño real del archivo
		remainingSize := int64(inode.I_size) - int64(content.Len())
		if remainingSize <= 0 {
			break
		}
		readSize := int64(len(fileBlock.B_content))
		if readSize > remainingSize {
			readSize = remainingSize
		}
		content.Write(fileBlock.B_content[:readSize])
	}
	return content.Bytes(), nil
}

// --- Funciones auxiliares de lectura de bajo nivel ---
func ReadInode(file *os.File, sb structs.Superblock, index int32) (structs.Inode, error) {
	var inode structs.Inode
	offset := int64(sb.S_inode_start) + int64(index)*int64(sb.S_inode_size)
	file.Seek(offset, 0)
	err := binary.Read(file, binary.BigEndian, &inode)
	return inode, err
}
func ReadFileBlock(file *os.File, sb structs.Superblock, index int32) (structs.FileBlock, error) {
	var block structs.FileBlock
	offset := int64(sb.S_block_start) + int64(index)*int64(sb.S_block_size)
	file.Seek(offset, 0)
	err := binary.Read(file, binary.BigEndian, &block)
	return block, err
}
func ReadFolderBlock(file *os.File, sb structs.Superblock, index int32) (structs.FolderBlock, error) {
	var block structs.FolderBlock
	offset := int64(sb.S_block_start) + int64(index)*int64(sb.S_block_size)
	file.Seek(offset, 0)
	err := binary.Read(file, binary.BigEndian, &block)
	return block, err
}

func ReadSuperblock(file *os.File, partitionStart int64) (structs.Superblock, error) {
	var sb structs.Superblock
	_, err := file.Seek(partitionStart, 0)
	if err != nil {
		return sb, fmt.Errorf("error al posicionar el cursor: %v", err)
	}
	err = binary.Read(file, binary.BigEndian, &sb)
	if err != nil {
		return sb, fmt.Errorf("error al leer el superbloque: %v", err)
	}
	return sb, nil
}

// CreateFile es la función principal para crear archivos.
func CreateFile(mountedPartition state.MountedPartition, path string, isRecursive bool, content []byte) error {
	file, err := os.OpenFile(mountedPartition.Path, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el disco: %w", err)
	}
	defer file.Close()

	sb, err := ReadSuperblock(file, mountedPartition.Start)
	if err != nil {
		return err
	}

	dirPath := filepath.Dir(path)
	fileName := filepath.Base(path)

	parentInode, parentInodeIndex, err := FindInodeByPath(file, sb, dirPath)
	if err != nil {
		if !isRecursive {
			return fmt.Errorf("la ruta padre '%s' no existe. Use el flag -r para crearla.", dirPath)
		}
		// --- INICIO DE LA LÓGICA RECURSIVA ---
		fmt.Printf("INFO: Creando carpetas recursivamente para la ruta '%s'\n", dirPath)
		parts := strings.Split(strings.Trim(dirPath, "/"), "/")
		currentPath := "/"
		var currentInodeIndex int32 = 0 // Empezamos desde el inodo raíz

		for _, part := range parts {
			if part == "" {
				continue
			}
			// Recargamos el superbloque en cada iteración por si cambió
			sb, _ = ReadSuperblock(file, mountedPartition.Start)

			// Intentamos encontrar la siguiente carpeta en la ruta
			nextInode, nextInodeIndex, findErr := FindInodeByPath(file, sb, filepath.Join(currentPath, part))
			if findErr != nil {
				// Si no existe, la creamos
				fmt.Printf("INFO: La carpeta '%s' no existe, creando...\n", part)
				// Creamos la nueva carpeta (tipo 0)
				newFolderInodeIndex, createErr := createNewFsObject(file, &sb, mountedPartition.Start, currentInodeIndex, part, 0)
				if createErr != nil {
					return fmt.Errorf("no se pudo crear la carpeta recursiva '%s': %w", part, createErr)
				}
				currentInodeIndex = newFolderInodeIndex
			} else {
				// Si existe, solo actualizamos el índice actual
				if nextInode.I_type != 0 {
					return fmt.Errorf("la ruta '%s' contiene un archivo, no una carpeta", filepath.Join(currentPath, part))
				}
				currentInodeIndex = nextInodeIndex
			}
			currentPath = filepath.Join(currentPath, part)
		}
		// Al final del bucle, currentInodeIndex es el índice del inodo de la carpeta padre final
		parentInodeIndex = currentInodeIndex
		parentInode, _ = ReadInode(file, sb, parentInodeIndex)
		// --- FIN DE LA LÓGICA RECURSIVA ---
	}

	_, err = findEntryInInode(file, sb, parentInode, fileName)

	if err == nil {
		return fmt.Errorf("el archivo o carpeta '%s' ya existe en '%s'", fileName, dirPath)
	}

	if !CheckPermissions(state.CurrentUser, parentInode, 2) { // 2 = Write
		return fmt.Errorf("permiso denegado: no tiene permisos de escritura en la carpeta padre")
	}

	newInodeNum, err := createNewFsObject(file, &sb, mountedPartition.Start, parentInodeIndex, fileName, 1) // 1 para archivo
	if err != nil {
		return err
	}

	if len(content) > 0 {
		// Recargar el superbloque por si cambió durante la creación del objeto
		sb, err = ReadSuperblock(file, mountedPartition.Start)
		if err != nil {
			return fmt.Errorf("error al recargar superbloque antes de escribir contenido: %w", err)
		}
		err = WriteContentToInode(file, &sb, mountedPartition.Start, newInodeNum, content)
		if err != nil {
			return fmt.Errorf("error al escribir contenido en el archivo: %w", err)
		}
		fmt.Printf("INFO: Se escribieron %d bytes en el archivo '%s'.\n", len(content), fileName)
	}
	// --- FIN DE LA CORRECCIÓN ---
	return nil
}

// createNewFsObject es una función auxiliar para crear inodos y bloques.
func createNewFsObject(file *os.File, sb *structs.Superblock, partStart int64, parentInodeNum int32, newName string, objectType byte) (int32, error) {
	// 0) Asegurar que el inodo padre tenga al menos un bloque disponible
	parentInode, err := ReadInode(file, *sb, parentInodeNum)
	if err != nil {
		return -1, err
	}
	// si el padre no tiene bloques asignados, asignar uno
	hasBlock := false
	for _, b := range parentInode.I_block {
		if b != -1 {
			hasBlock = true
			break
		}
	}
	if !hasBlock {
		parentNewBlock, err := utils.FindFreeBit(file, int64(sb.S_bm_block_start), sb.S_blocks_count)
		if err != nil {
			return -1, err
		}
		// inicializar bloque de carpeta vacío
		fb := structs.FolderBlock{}
		for i := range fb.B_content {
			fb.B_content[i].B_inodo = -1
		}
		blockPos := int64(sb.S_block_start) + int64(parentNewBlock)*int64(sb.S_block_size)
		if err := utils.WriteObject(file, blockPos, fb); err != nil {
			return -1, err
		}
		// asignar al primer slot libre del padre
		for i := range parentInode.I_block {
			if parentInode.I_block[i] == -1 {
				parentInode.I_block[i] = parentNewBlock
				break
			}
		}
		// marcar bitmap y actualizar superbloque
		if err := utils.SetBit(file, int64(sb.S_bm_block_start), parentNewBlock, '1'); err != nil {
			return -1, err
		}
		sb.S_free_blocks_count--
		// escribir inodo padre actualizado
		parentPos := int64(sb.S_inode_start) + int64(parentInodeNum)*int64(sb.S_inode_size)
		if err := utils.WriteObject(file, parentPos, parentInode); err != nil {
			return -1, err
		}
	}

	// 1) Buscar inodo y bloque libres para el nuevo objeto
	newInodeIndex, err := utils.FindFreeBit(file, int64(sb.S_bm_inode_start), sb.S_inodes_count)
	if err != nil {
		return -1, err
	}
	newBlockIndex, err := utils.FindFreeBit(file, int64(sb.S_bm_block_start), sb.S_blocks_count)
	if err != nil {
		return -1, err
	}

	// 2) Crear y escribir el nuevo inodo
	newInode := structs.Inode{
		I_uid:   state.CurrentUser.Id,
		I_gid:   state.CurrentUser.Group,
		I_size:  0,
		I_atime: time.Now().Unix(),
		I_ctime: time.Now().Unix(),
		I_mtime: time.Now().Unix(),
		I_type:  int32(objectType),
		I_perm:  664,
	}
	for i := range newInode.I_block {
		newInode.I_block[i] = -1
	}
	newInode.I_block[0] = newBlockIndex

	inodePosition := int64(sb.S_inode_start) + int64(newInodeIndex)*int64(sb.S_inode_size)
	if err := utils.WriteObject(file, inodePosition, newInode); err != nil {
		return -1, err
	}

	// 3) Inicializar el bloque del nuevo objeto
	blockPosition := int64(sb.S_block_start) + int64(newBlockIndex)*int64(sb.S_block_size)
	if objectType == 1 { // Archivo
		newFileBlock := structs.FileBlock{}
		if err := utils.WriteObject(file, blockPosition, newFileBlock); err != nil {
			return -1, err
		}
	} else { // Carpeta
		newFolderBlock := structs.FolderBlock{}
		for i := range newFolderBlock.B_content {
			newFolderBlock.B_content[i].B_inodo = -1
		}
		if err := utils.WriteObject(file, blockPosition, newFolderBlock); err != nil {
			return -1, err
		}
	}

	// 4) Añadir la entrada al directorio padre
	if err := addEntryToFolder(file, sb, parentInodeNum, newName, newInodeIndex); err != nil {
		return -1, err
	}

	// 5) Marcar bits y actualizar superbloque
	if err := utils.SetBit(file, int64(sb.S_bm_inode_start), newInodeIndex, '1'); err != nil {
		return -1, err
	}
	if err := utils.SetBit(file, int64(sb.S_bm_block_start), newBlockIndex, '1'); err != nil {
		return -1, err
	}
	sb.S_free_inodes_count--
	sb.S_free_blocks_count--
	sb.S_first_ino, _ = utils.FindFreeBit(file, int64(sb.S_bm_inode_start), sb.S_inodes_count)
	sb.S_first_blo, _ = utils.FindFreeBit(file, int64(sb.S_bm_block_start), sb.S_blocks_count)

	// 6) Escribir superbloque actualizado
	if err := utils.WriteObject(file, partStart, *sb); err != nil {
		return -1, err
	}
	return newInodeIndex, nil
}

// WriteContentToInode escribe un slice de bytes en los bloques de un inodo.
func WriteContentToInode(file *os.File, sb *structs.Superblock, partStart int64, inodeNum int32, content []byte) error {
	inode, err := ReadInode(file, *sb, inodeNum)
	if err != nil {
		return err
	}

	bytesWritten := 0
	blockPointerIndex := 0

	for bytesWritten < len(content) {
		if blockPointerIndex >= 12 {
			return fmt.Errorf("archivo demasiado grande, no se implementan punteros indirectos")
		}

		blockNum := inode.I_block[blockPointerIndex]
		if blockNum == -1 {
			newBlock, err := utils.FindFreeBit(file, int64(sb.S_bm_block_start), sb.S_blocks_count)
			if err != nil {
				return err
			}
			blockNum = newBlock
			inode.I_block[blockPointerIndex] = blockNum
			if err := utils.SetBit(file, int64(sb.S_bm_block_start), blockNum, '1'); err != nil {
				return err
			}
			sb.S_free_blocks_count--
		}

		blockPosition := int64(sb.S_block_start) + int64(blockNum)*int64(sb.S_block_size)
		fileBlock := structs.FileBlock{}

		remainingContent := content[bytesWritten:]
		chunkSize := len(remainingContent)
		if chunkSize > int(sb.S_block_size) {
			chunkSize = int(sb.S_block_size)
		}
		copy(fileBlock.B_content[:], remainingContent[:chunkSize])

		if err := utils.WriteObject(file, blockPosition, fileBlock); err != nil {
			return err
		}

		bytesWritten += chunkSize
		blockPointerIndex++
	}
	// Actualizar el tamaño del inodo con la cantidad de bytes del contenido
	inode.I_size = int32(len(content))
	inode.I_mtime = time.Now().Unix()
	inodePosition := int64(sb.S_inode_start) + int64(inodeNum)*int64(sb.S_inode_size)
	if err := utils.WriteObject(file, inodePosition, inode); err != nil {
		return err
	}

	// Actualizar primer bloque libre y escribir superbloque
	if nextFree, _ := utils.FindFreeBit(file, int64(sb.S_bm_block_start), sb.S_blocks_count); nextFree >= 0 {
		sb.S_first_blo = nextFree
	}
	return utils.WriteObject(file, partStart, *sb)
}

// addEntryToFolder añade una nueva entrada a un bloque de carpeta.
func addEntryToFolder(file *os.File, sb *structs.Superblock, parentInodeNum int32, newName string, newInodeNum int32) error {
	parentInode, err := ReadInode(file, *sb, parentInodeNum)
	if err != nil {
		return err
	}

	// Buscar espacio en bloques existentes
	for _, blockPtr := range parentInode.I_block {
		if blockPtr == -1 {
			continue
		}
		folderBlock, err := ReadFolderBlock(file, *sb, blockPtr)
		if err != nil {
			return err
		}
		for i := range folderBlock.B_content {
			if folderBlock.B_content[i].B_inodo == -1 {
				copy(folderBlock.B_content[i].B_name[:], newName)
				folderBlock.B_content[i].B_inodo = newInodeNum
				blockPosition := int64(sb.S_block_start) + int64(blockPtr)*int64(sb.S_block_size)
				return utils.WriteObject(file, blockPosition, folderBlock)
			}
		}
	}

	// Si no hay espacio, buscar un slot libre en el inodo para asignar un nuevo bloque
	freeSlotIndex := -1
	for i, blockPtr := range parentInode.I_block {
		if blockPtr == -1 {
			freeSlotIndex = i
			break
		}
	}

	if freeSlotIndex == -1 {
		return errors.New("el directorio padre ha alcanzado el límite máximo de bloques (12)")
	}

	// Asignar un nuevo bloque para el directorio
	newBlockIndex, err := utils.FindFreeBit(file, int64(sb.S_bm_block_start), sb.S_blocks_count)
	if err != nil {
		return fmt.Errorf("no hay bloques libres disponibles: %w", err)
	}

	// Inicializar el nuevo bloque de carpeta
	newFolderBlock := structs.FolderBlock{}
	for i := range newFolderBlock.B_content {
		newFolderBlock.B_content[i].B_inodo = -1
	}

	// Agregar la nueva entrada al primer slot del nuevo bloque
	copy(newFolderBlock.B_content[0].B_name[:], newName)
	newFolderBlock.B_content[0].B_inodo = newInodeNum

	// Escribir el nuevo bloque
	blockPosition := int64(sb.S_block_start) + int64(newBlockIndex)*int64(sb.S_block_size)
	if err := utils.WriteObject(file, blockPosition, newFolderBlock); err != nil {
		return err
	}

	// Marcar el bloque como usado en el bitmap
	if err := utils.SetBit(file, int64(sb.S_bm_block_start), newBlockIndex, '1'); err != nil {
		return err
	}

	// Actualizar el inodo del directorio padre para referenciar el nuevo bloque
	parentInode.I_block[freeSlotIndex] = newBlockIndex
	parentInode.I_size += int32(sb.S_block_size)
	inodePosition := int64(sb.S_inode_start) + int64(parentInodeNum)*int64(sb.S_inode_size)
	if err := utils.WriteObject(file, inodePosition, parentInode); err != nil {
		return err
	}

	// Actualizar el superbloque
	sb.S_free_blocks_count--
	if nextFree, _ := utils.FindFreeBit(file, int64(sb.S_bm_block_start), sb.S_blocks_count); nextFree >= 0 {
		sb.S_first_blo = nextFree
	}

	return nil
}

// findEntryInInode busca una entrada por nombre dentro de los bloques de un inodo de carpeta.
func findEntryInInode(file *os.File, sb structs.Superblock, inode structs.Inode, entryName string) (int32, error) {
	for _, blockPtr := range inode.I_block {
		if blockPtr == -1 {
			continue
		}
		folderBlock, err := ReadFolderBlock(file, sb, blockPtr)
		if err != nil {
			return -1, err
		}
		for _, entry := range folderBlock.B_content {
			if entry.B_inodo != -1 && strings.TrimRight(string(entry.B_name[:]), "\x00") == entryName {
				return entry.B_inodo, nil
			}
		}
	}
	return -1, errors.New("entrada no encontrada")
}

// CheckPermissions verifica si un usuario tiene los permisos requeridos sobre un inodo.
func CheckPermissions(user state.User, inode structs.Inode, requiredPerm int32) bool {
	// El usuario root (uid=1) siempre tiene todos los permisos.
	if user.Id == 1 {
		return true
	}

	// Extraer permisos para propietario, grupo y otros.
	ownerPerm := (inode.I_perm >> 6) & 7 // Permisos del propietario
	groupPerm := (inode.I_perm >> 3) & 7 // Permisos del grupo
	otherPerm := inode.I_perm & 7        // Permisos para otros

	// Comprobar si el usuario es el propietario del inodo.
	if inode.I_uid == user.Id {
		return (ownerPerm & requiredPerm) == requiredPerm
	}

	// Comprobar si el usuario pertenece al grupo del inodo.
	if inode.I_gid == user.Group {
		return (groupPerm & requiredPerm) == requiredPerm
	}

	// Si no es propietario ni del grupo, aplicar permisos de "otros".
	return (otherPerm & requiredPerm) == requiredPerm
}

// CreateDirectory crea un nuevo directorio en el sistema de archivos
func CreateDirectory(mountedPartition state.MountedPartition, path string, isRecursive bool) error {
	file, err := os.OpenFile(mountedPartition.Path, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el disco: %w", err)
	}
	defer file.Close()

	sb, err := ReadSuperblock(file, mountedPartition.Start)
	if err != nil {
		return err
	}

	dirPath := filepath.Dir(path)
	dirName := filepath.Base(path)

	parentInode, parentInodeIndex, err := FindInodeByPath(file, sb, dirPath)
	if err != nil {
		if !isRecursive {
			return fmt.Errorf("la ruta padre '%s' no existe. Use el flag -p para crearla", dirPath)
		}
		// Crear directorios padre recursivamente
		err = CreateDirectory(mountedPartition, dirPath, true)
		if err != nil {
			return fmt.Errorf("error creando directorios padre: %w", err)
		}
		// Recargamos después de crear los padres
		sb, err = ReadSuperblock(file, mountedPartition.Start)
		if err != nil {
			return err
		}
		parentInode, parentInodeIndex, err = FindInodeByPath(file, sb, dirPath)
		if err != nil {
			return fmt.Errorf("error encontrando padre después de crearlo: %w", err)
		}
	}

	// Verificar que no existe ya
	_, err = findEntryInInode(file, sb, parentInode, dirName)
	if err == nil {
		return fmt.Errorf("el directorio '%s' ya existe en '%s'", dirName, dirPath)
	}

	// Verificar permisos
	if !CheckPermissions(state.CurrentUser, parentInode, 2) { // 2 = Write
		return fmt.Errorf("permiso denegado: no tiene permisos de escritura en la carpeta padre")
	}

	// Crear el nuevo directorio (tipo 0)
	_, err = createNewFsObject(file, &sb, mountedPartition.Start, parentInodeIndex, dirName, 0)
	if err != nil {
		return fmt.Errorf("error creando directorio: %w", err)
	}

	return nil
}
