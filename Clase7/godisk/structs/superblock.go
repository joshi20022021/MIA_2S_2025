// structs/superblock.go
package structs

type Superblock struct {
	S_filesystem_type   int32 // Tipo de sistema de archivos (2 para 2fs)
	S_inodes_count      int32 // Número total de inodos
	S_blocks_count      int32 // Número total de bloques
	S_free_blocks_count int32 // Número de bloques libres
	S_free_inodes_count int32 // Número de inodos libres
	S_mtime             int64 // Última fecha de montaje (Unix timestamp)
	S_umtime            int64 // Última fecha de desmontaje (Unix timestamp)
	S_mnt_count         int32 // Contador de montajes
	S_magic             int32 // Valor mágico: 0xEF53
	S_inode_size        int32 // Tamaño del inodo (sizeof)
	S_block_size        int32 // Tamaño del bloque (sizeof)
	S_first_ino         int32 // Primer inodo libre (apuntador)
	S_first_blo         int32 // Primer bloque libre (apuntador)
	S_bm_inode_start    int32 // Inicio del bitmap de inodos
	S_bm_block_start    int32 // Inicio del bitmap de bloques
	S_inode_start       int32 // Inicio de la tabla de inodos
	S_block_start       int32 // Inicio de la tabla de bloques
}
