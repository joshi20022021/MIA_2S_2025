# Comandos Básicos
# MKDISK
### Crear un disco de 10 MB
```bash
mkdisk -size=10 -path=/home/joshipanda2002/Discos/disco1.mia
```
### Crear un disco de 3000 KB (aproximadamente 2.93 MB)
```bash
mkdisk -size=3000 -unit=K -path=/home/joshipanda2002/Discos/disco2.mia
```

## Crear un disco de 5 MB con ajuste "Best Fit"
```bash
mkdisk -size=5 -fit=BF -path=/home/joshipanda2002/Discos/disco3.mia
```

## Crear un disco de 15 MB en el directorio actual
```bash
mkdisk -size=15 -path=disco5
```

## Crear automáticamente subdirectorios subdir1/subdir2 y un disco de 20 MB
```bash
mkdisk -size=20 -path=/home/joshipanda2002/Discos/subdir1/subdir2/disco6.mia
```

# FDISK

## Crear una partición primaria de 15 MB con el peor ajuste (Worst Fit - wf)
```bash
fdisk -size=15 -unit=m -path=./mi_disco.mia -name=Primaria1 -type=p -fit=wf
```

## Crear otra partición primaria de 5120 Kilobytes con el mejor ajuste (Best Fit - bf)
```bash
fdisk -size=5120 -unit=k -path=./mi_disco.mia -name=Primaria2 -type=p -fit=bf
```

## Crear una partición extendida de 50 MB con el primer ajuste (First Fit - ff)
```bash
fdisk -size=50 -unit=m -path=./mi_disco.mia -name=Extendida -type=e -fit=ff
```

## Crear una partición lógica de 10 MB
```bash
fdisk -size=10 -unit=m -path=./mi_disco.mia -name=LogicaA -type=l
```

## Crear otra partición lógica de 20480 bytes
```bash
fdisk -size=20480 -unit=b -path=./mi_disco.mia -name=LogicaB -type=l
```