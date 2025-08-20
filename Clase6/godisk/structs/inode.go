// structs/inode.go
package structs

type Inode struct {
	I_uid   int32     // UID del propietario
	I_gid   int32     // GID del grupo
	I_size  int32     // Tamaño del archivo en bytes
	I_atime int64     // Última fecha de acceso (Unix timestamp)
	I_ctime int64     // Fecha de creación (Unix timestamp)
	I_mtime int64     // Última fecha de modificación (Unix timestamp)
	I_block [15]int32 // Bloques de datos (12 directos, 1 indirecto, 1 doble, 1 triple)
	I_type  int32     // Tipo (1: archivo, 0: carpeta)
	I_perm  int32     // Permisos (formato UGO - ej. 664)
}
