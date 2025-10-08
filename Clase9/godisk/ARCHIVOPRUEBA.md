#Calificacion Proyecto 1
#2S 2025
#Cambiar "su_usuario" -> por el usuario de su distribución de linux  
#Cambiar “34” -> por los ultimos dos digitos de su carnet  
#----------------- 1. MKDISK  -----------------
#----------------- MKDISK CON ERROR -----------------
# ERROR PARAMETROS
mkdisk -param=x -size=30 -path=/home/”su_usuario”/Calificacion_MIA/Discos/DiscoN.mia
#----------------- CREACION DE DISCOS -----------------
# ERROR PARAMETROS
mkdisk -tamaño=3000 -unit=K -path=/home/joshipanda2002/Calificacion_MIA/Discos/DiscoN.mia
# 50M A
Mkdisk -size=50 -unit=M -fit=FF -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco1.mia
# 50M B
Mkdisk -unit=k -size=51200 -fit=BF -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco2.mia
# 13M C
mkDisk -size=13 -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco3.mia
# 50M D
mkdisk -size=51200 -unit=K -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco4.mia
# 20M E
mkDisk -size=20 -unit=M -fit=WF -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco5.mia
# 50M F X
Mkdisk -size=50 -unit=M -fit=FF -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco6.mia
# 50M G X
Mkdisk -size=50 -unit=M -fit=FF -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco7.mia
# 50M H X
mkdisk -size=51200 -unit=K -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco8.mia
# 50M I X
mkdisk -size=51200 -unit=K -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco9.mia
# 50M J X
mkdisk -size=51200 -unit=K -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco10.mia
#-----------------2. RMDISK-----------------
#ERROR DISCO NO EXISTE
rmdisk -path=/home/joshipanda2002/Calificacion_MIA/Discos/DiscoN.mia
# BORRANDO DISCO
rmdisk -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco6.mia
# BORRANDO DISCO
rmdisk -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco7.mia
# BORRANDO DISCO
rmdisk -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco8.mia
# BORRANDO DISCO
rmdisk -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco9.mia
# BORRANDO DISCO
rmdisk -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco10.mia
#-----------------3. FDISK-----------------
#-----------------CREACION DE PARTICIONES-----------------
#DISCO 1
# ERROR RUTA NO ENCONTRADA
fdisk -type=P -unit=b -name=PartErr -size=10485760 -path=/home/joshipanda2002/Calificacion_MIA/Discos/DiscoN.mia -fit=BF 
# PRIMARIA 10M
fdisk -type=P -unit=b -name=Part11 -size=10485760 -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco1.mia -fit=BF
# PRIMARIA 10M
fdisk -type=P -unit=k -name=Part12 -size=10240 -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco1.mia -fit=BF
# PRIMARIA 10M
fdisk -type=P -unit=M -name=Part13 -size=10 -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco1.mia -fit=BF
# PRIMARIA 10M
fdisk -type=P -unit=b -name=Part14 -size=10485760 -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco1.mia -fit=BF
#ERR LMITE PARTICION PRIMARIA
fdisk -type=P -unit=b -name=PartErr -size=10485760 -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco1.mia -fit=BF
# LIBRE DISCO 1: 50-4*10 = 10 -> 20%
#DISCO 3
# ERROR FALTA ESPACIO
fdisk -type=P -unit=m -name=PartErr -size=20 -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco3.mia
#4M
fdisk -type=P -unit=m -name=Part31 -size=4 -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco3.mia
#4M
fdisk -type=P -unit=m -name=Part32 -size=4 -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco3.mia
#1M
fdisk -type=P -unit=m -name=Part33 -size=1 -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco3.mia
#LIBRE DISCO 3: 13-9= 4 -> 30.77%
#DISCO 5
# 5MB
fdisk -type=E -unit=k -name=Part51 -size=5120 -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco5.mia -fit=BF
# 1MB
fdisk -type=L -unit=k -name=Part52 -size=1024 -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco5.mia -fit=BF
# 5MB
fdisk -type=P -unit=k -name=Part53 -size=5120 -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco5.mia -fit=BF
# 1MB
fdisk -type=L -unit=k -name=Part54 -size=1024 -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco5.mia -fit=BF
# 1MB
fdisk -type=L -unit=k -name=Part55 -size=1024 -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco5.mia -fit=BF
# 1MB
fdisk -type=L -unit=k -name=Part56 -size=1024 -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco5.mia -fit=BF
# LIBRE DISCO 5: 20-10 = 5 -> 50%
# LIBRE EXTENDIDA 2: 5-4 = 1 -> 20% (por los EBR deberia ser menos)
#-----------------MOUNT-----------------
#-----------------MONTAR PARTICIONES-----------------
#DISCO 1
#341A -> 34 sus ultimos dos digitos del carnet
mount -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco1.mia -name=Part11
#342A -> 34 sus ultimos dos digitos del carnet
mount -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco1.mia -name=Part12
#ERROR PARTICION YA MONTADA
mount -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco1.mia -name=Part11
#DISCO 3
#ERROR PARTCION NO EXISTE
mount -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco3.mia -name=Part0
#341B -> 34 sus ultimos dos digitos del carnet
mount -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco3.mia -name=Part31
#342B -> 34 sus ultimos dos digitos del carnet
mount -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco3.mia -name=Part32
#DISCO 5
#341C -> 34 sus ultimos dos digitos del carnet
mount -path=/home/joshipanda2002/Calificacion_MIA/Discos/Disco5.mia -name=Part53
#-----------------MOUNT-----------------
# ENLISTAR
mounted
#-----------------REPORTES PARTE 1-----------------
#DISCO 1
#ERROR ID NO ENCONTRADO
rep -id=A341 -path=/home/joshipanda2002/Calificacion_MIA/Reportes/p1_rE.jpg -name=mbr
#REPORTE DISK
rep -id=341A -path=/home/joshipanda2002/Calificacion_MIA/Reportes/p1_r1_disk.jpg -name=disk
#REPORTE MBR 
rep -id=341A -path=/home/joshipanda2002/Calificacion_MIA/Reportes/p1_r2_mbr.jpg -name=mbr
#DISCO 3
#ERROR ID NO ENCONTRADO
rep -id=343B -path=/home/joshipanda2002/Calificacion_MIA/Reportes/p1_rE_mbr.jpg -name=mbr
#REPORTE DISK
rep -id=341B -path=/home/joshipanda2002/Calificacion_MIA/Reportes/p1_r3_disk.jpg -name=disk
#REPORTE MBR
rep -id=342B -path=/home/joshipanda2002/Calificacion_MIA/Reportes/p1_r4_disk.jpg -name=mbr
#DISCO 5
#ERROR ID NO ENCONTRADO
rep -id=IDx -path=/home/joshipanda2002/Calificacion_MIA/Reportes/p1_rE_mbr.jpg -name=mbr
#REPORTE DISK
rep -id=341C -path=/home/joshipanda2002/Calificacion_MIA/Reportes/p1_r5_disk.jpg -name=disk
#REPORTE MBR
rep -id=341C -path=/home/joshipanda2002/Calificacion_MIA/Reportes/p1_r6_mbr.jpg -name=mbr
#-----------------5. MKFS-----------------
mkfs -type=full -id=341A
#-----------------PARTE 3-----------------
#-----------------7. LOGIN-----------------
login -user=root -pass=123 -id=341A
#ERROR SESION INICIADA
login -user=root -pass=123 -id=341A
#-----------------14. MKFILE-----------------
mkfile -r -path=/home/archivos/user/docs/Tarea.txt -size=75 
mkfile -r -path=/home/archivos/user/docs/Tarea2.txt -size=768
#Para este comando hay que crear un archivo en la computadora y en cont poner su primer nombre
#Crear un archivo txt en su Escritorio llamado NAME
# Cambiar la ruta del cont por la del archivo NAME.txt que creo
mkfile -r -path=/home/archivos/user/docs/Tarea3.txt -size=10 -jose=/home/”su_usuario”/Calificacion_MIA/CONT/NAME.txt
#ERROR NO EXISTE RUTA
mkfile -path="/home/archivos/noexiste/b1.txt"
#ERROR NEGATIVO
mkfile -path="/home/archivos/b1.txt" -size=-45
#RECURSIVO
mkfile -r -path=/home/archivos/user/docs/usac/archivos/proyectos/fase1/entrada.txt
#-----------------JOURNALING-----------------------------
journal -id=341A
#-----------------6. CAT-----------------
cat -file1=/home/archivos/user/docs/Tarea2.txt
cat -file1=/home/archivos/user/docs/Tarea3.txt
#------------------------8. LOGOUT------------------------
logout