package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"godisk/structs"
	"math"
	"os"
	"os/exec"
	"reflect"
	"sort"
	"strings"
)

// FreeSpace es una estructura exportada para manejar los espacios libres.
type FreeSpace struct {
	Start int64
	End   int64
	Size  int64
}

// GetFreeSpaces analiza el MBR y devuelve una lista de todos los huecos libres.
func GetFreeSpaces(mbr *structs.MBR) []FreeSpace {
	var spaces []FreeSpace
	var occupiedPartitions []structs.Partition

	for _, p := range mbr.Mbr_partitions {
		if p.Part_status == '1' {
			occupiedPartitions = append(occupiedPartitions, p)
		}
	}

	sort.Slice(occupiedPartitions, func(i, j int) bool {
		return occupiedPartitions[i].Part_start < occupiedPartitions[j].Part_start
	})

	currentPos := int64(binary.Size(structs.MBR{}))

	for _, p := range occupiedPartitions {
		if p.Part_start > currentPos {
			size := p.Part_start - currentPos
			spaces = append(spaces, FreeSpace{Start: currentPos, End: p.Part_start - 1, Size: size})
		}
		currentPos = p.Part_start + p.Part_s
	}

	if currentPos < mbr.Mbr_tamano {
		size := mbr.Mbr_tamano - currentPos
		spaces = append(spaces, FreeSpace{Start: currentPos, End: mbr.Mbr_tamano - 1, Size: size})
	}

	return spaces
}

// GetFreeSpacesInExtended analiza los EBRs y devuelve los huecos en una extendida.
func GetFreeSpacesInExtended(extended structs.Partition, logicals []structs.EBR) []FreeSpace {
	var spaces []FreeSpace

	sort.Slice(logicals, func(i, j int) bool {
		return logicals[i].Part_start < logicals[j].Part_start
	})

	currentPos := extended.Part_start
	for _, l := range logicals {
		ebrStart := l.Part_start - int64(binary.Size(l))
		if ebrStart > currentPos {
			size := ebrStart - currentPos
			spaces = append(spaces, FreeSpace{Start: currentPos, End: ebrStart - 1, Size: size})
		}
		currentPos = l.Part_start + l.Part_s
	}

	if currentPos < (extended.Part_start + extended.Part_s) {
		size := (extended.Part_start + extended.Part_s) - currentPos
		spaces = append(spaces, FreeSpace{Start: currentPos, End: (extended.Part_start + extended.Part_s) - 1, Size: size})
	}

	return spaces
}

// FindFirstFit encuentra el primer hueco que sea lo suficientemente grande.
func FindFirstFit(spaces []FreeSpace, requiredSize int64) int64 {
	for _, space := range spaces {
		if space.Size >= requiredSize {
			return space.Start
		}
	}
	return -1
}

// FindBestFit encuentra el hueco que deje el menor desperdicio.
func FindBestFit(spaces []FreeSpace, requiredSize int64) int64 {
	bestStart := int64(-1)
	minDifference := int64(math.MaxInt64)

	for _, space := range spaces {
		if space.Size >= requiredSize {
			difference := space.Size - requiredSize
			if difference < minDifference {
				minDifference = difference
				bestStart = space.Start
			}
		}
	}
	return bestStart
}

// FindWorstFit encuentra el hueco más grande disponible.
func FindWorstFit(spaces []FreeSpace, requiredSize int64) int64 {
	worstStart := int64(-1)
	maxSize := int64(-1)

	for _, space := range spaces {
		if space.Size >= requiredSize {
			if space.Size > maxSize {
				maxSize = space.Size
				worstStart = space.Start
			}
		}
	}
	return worstStart
}

// WriteMBR, ReadMBR, WriteEBR, ReadEBR...
func WriteMBR(file *os.File, mbr *structs.MBR) error {
	file.Seek(0, 0)
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.LittleEndian, mbr)
	if err != nil {
		return fmt.Errorf("error al serializar el MBR: %w", err)
	}
	_, err = file.Write(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("error al escribir el MBR en el disco: %w", err)
	}
	return nil
}

func ReadMBR(file *os.File) (structs.MBR, error) {
	var mbr structs.MBR
	file.Seek(0, 0)
	err := binary.Read(file, binary.LittleEndian, &mbr)
	return mbr, err
}

func WriteEBR(file *os.File, ebr *structs.EBR, start int64) error {
	file.Seek(start, 0)
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.LittleEndian, ebr)
	if err != nil {
		return fmt.Errorf("error al serializar el EBR: %w", err)
	}
	_, err = file.Write(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("error al escribir el EBR en el disco: %w", err)
	}
	return nil
}

func ReadEBR(file *os.File, start int64) (structs.EBR, error) {
	var ebr structs.EBR
	file.Seek(start, 0)
	err := binary.Read(file, binary.LittleEndian, &ebr)
	if err != nil {
		return structs.EBR{}, fmt.Errorf("error al leer el EBR: %w", err)
	}
	return ebr, nil
}

// Sizeof calcula el tamaño en bytes de una estructura.
func Sizeof(v interface{}) int {
	return int(reflect.TypeOf(v).Size())
}
func ReadSuperblock(file *os.File, partitionStart int64) (structs.Superblock, error) {
	var sb structs.Superblock
	file.Seek(partitionStart, 0)
	err := binary.Read(file, binary.BigEndian, &sb)
	return sb, err
}

func ReadBytes(file *os.File, start int64, size int64) ([]byte, error) {
	bytes := make([]byte, size)
	_, err := file.ReadAt(bytes, start)
	if err != nil {
		return nil, fmt.Errorf("error al leer bytes: %w", err)
	}
	return bytes, nil
}

func ReadInode(file *os.File, start int64) (structs.Inode, error) {
	var inode structs.Inode
	size := int64(binary.Size(inode))
	data, err := ReadBytes(file, start, size)
	if err != nil {
		return structs.Inode{}, err
	}
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &inode)
	if err != nil {
		return structs.Inode{}, fmt.Errorf("error al decodificar el inodo: %w", err)
	}
	return inode, nil
}

func GenerateReport(dotContent string, outputPath string) error {
	cmd := exec.Command("dot", "-Tjpg", "-o", outputPath)
	cmd.Stdin = strings.NewReader(dotContent)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error al ejecutar dot: %v\n%s", err, stderr.String())
	}
	return nil
}

func ReadFolderBlock(file *os.File, start int64) (structs.FolderBlock, error) {
	var block structs.FolderBlock
	size := int64(binary.Size(block))
	data, err := ReadBytes(file, start, size)
	if err != nil {
		return structs.FolderBlock{}, err
	}
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &block)
	if err != nil {
		return structs.FolderBlock{}, fmt.Errorf("error al decodificar el bloque de carpeta: %w", err)
	}
	return block, nil
}
