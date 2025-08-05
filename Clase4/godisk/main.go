package main

// Importaciones necesarias para el funcionamiento del programa
import (
	"bufio"           // Paquete para leer la entrada de usuario línea por línea de manera eficiente
	"flag"            // Paquete para manejar flags/parámetros de línea de comandos
	"fmt"             // Paquete para formatear e imprimir texto en consola
	"godisk/commands" // Importa nuestro paquete personalizado que contiene los comandos disponibles
	"os"              // Paquete para interactuar con el sistema operativo (stdin, stderr, etc.)
	"strings"         // Paquete para manipular y modificar cadenas de texto
)

func main() {
	// bufio.Scanner lee línea por línea
	scanner := bufio.NewScanner(os.Stdin)

	// Bucle infinito para mantener el programa en ejecución hasta que el usuario decida salir
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break // Si no puede leer (EOF o error), sale del bucle
		}
		// Obtiene el texto de la línea leída, sin el carácter de nueva línea
		line := scanner.Text()
		// Verifica si el usuario quiere salir del programa
		if strings.ToLower(line) == "exit" {
			fmt.Println("Saliendo...")
			break // Termina el bucle y por tanto el programa
		}
		// strings.Fields() es más robusto que strings.Split() porque maneja múltiples espacios
		parts := strings.Fields(line)
		// Si el usuario solo presionó Enter (línea vacía), no hay nada que procesar
		if len(parts) == 0 {
			continue // Vuelve al inicio del bucle para mostrar el prompt nuevamente
		}

		// La primera palabra siempre es el comando a ejecutar
		command := strings.ToLower(parts[0]) // Convertir a minúsculas para comparación
		// El resto de las palabras son los argumentos y flags del comando
		args := parts[1:] // Slice que incluye desde el índice 1 hasta el final

		// Switch para determinar qué comando ejecutar
		switch command {
		case "mkdisk":
			mkdiskCmd := flag.NewFlagSet("mkdisk", flag.ExitOnError)
			size := mkdiskCmd.Int("size", 0, "Tamaño del disco.")
			unit := mkdiskCmd.String("unit", "m", "Unidad del tamaño (k/m).")
			fit := mkdiskCmd.String("fit", "ff", "Tipo de ajuste (bf/ff/wf).")
			path := mkdiskCmd.String("path", "", "Ruta del disco a crear.")
			mkdiskCmd.Parse(args)

			if *path == "" {
				fmt.Println("Error: el parámetro -path es obligatorio para mkdisk.")
				continue
			}
			if *size <= 0 {
				fmt.Println("Error: el parámetro -size es obligatorio y debe ser positivo.")
				continue
			}
			commands.ExecuteMkdisk(*size, *unit, *fit, *path)

		case "rmdisk":
			rmdiskCmd := flag.NewFlagSet("rmdisk", flag.ExitOnError)
			path := rmdiskCmd.String("path", "", "Ruta del disco a eliminar.")
			rmdiskCmd.Parse(args)

			if *path == "" {
				fmt.Println("Error: el parámetro -path es obligatorio para rmdisk.")
				continue
			}
			commands.ExecuteRmdisk(*path)
		case "fdisk":
			fdiskCmd := flag.NewFlagSet("fdisk", flag.ExitOnError)
			size := fdiskCmd.Int64("size", 0, "Tamaño de la partición.")
			path := fdiskCmd.String("path", "", "Ruta del disco.")
			name := fdiskCmd.String("name", "", "Nombre de la partición.")
			unit := fdiskCmd.String("unit", "k", "Unidad del tamaño (b/k/m).")
			typeStr := fdiskCmd.String("type", "p", "Tipo de partición (p/e/l).")
			fit := fdiskCmd.String("fit", "wf", "Tipo de ajuste (bf/ff/wf).")
			fdiskCmd.Parse(args)

			if *path == "" || *name == "" || *size <= 0 {
				fmt.Println("Error: los parámetros -path, -name y -size son obligatorios para fdisk.")
				continue
			}
			commands.ExecuteFdisk(*path, *name, *unit, *typeStr, *fit, *size)

		default:
			fmt.Printf("Comando '%s' no reconocido.\n", command)
		}
	}

	// Manejo de errores del scanner
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error leyendo la entrada:", err)
	}
}
