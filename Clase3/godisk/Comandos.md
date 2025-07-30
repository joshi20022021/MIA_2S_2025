# Comandos Básicos

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