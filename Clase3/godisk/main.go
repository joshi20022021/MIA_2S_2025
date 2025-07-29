package main

import (
	"bufio" // Paquete para leer la entrada
	"flag"
	"fmt"
	"godisk/commands"
	"os"
	"strings" // Paquete modificar texto
)

func main() {
	// Crea un nuevo scanner
	scanner := bufio.NewScanner(os.Stdin)

	// mantener programa en ejecución
	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		// Obtiene la línea
		line := scanner.Text()

		// salir del programa.
		if strings.ToLower(line) == "exit" {
			fmt.Println("Saliendo...")
			break
		}

		// crea un slice de strings
		parts := strings.Fields(line)

		// vuelve a esperar la entrada
		if len(parts) == 0 {
			continue
		}

		// La primera palabra es el comando.
		command := strings.ToLower(parts[0])
		// El resto de las palabras son los argumentos y flags.
		args := parts[1:]

		// partes del comando mkdisk
		switch command {
		case "mkdisk":
			mkdiskCmd := flag.NewFlagSet("mkdisk", flag.ExitOnError)
			size := mkdiskCmd.Int("size", 0, "Tamaño del disco.")
			unit := mkdiskCmd.String("unit", "m", "Unidad del tamaño (k/m).")
			fit := mkdiskCmd.String("fit", "ff", "Tipo de ajuste (bf/ff/wf).")
			path := mkdiskCmd.String("path", "", "Ruta del disco a crear.")

			// Parsea los argumentos
			mkdiskCmd.Parse(args)

			// validaciones
			if *path == "" {
				fmt.Println("Error: el parámetro -path es obligatorio para mkdisk.")
				continue // Vuelve al inicio del bucle.
			}
			if *size <= 0 {
				fmt.Println("Error: el parámetro -size es obligatorio y debe ser positivo.")
				continue
			}
			commands.ExecuteMkdisk(*size, *unit, *fit, *path)

		default:
			fmt.Printf("Comando '%s' no reconocido.\n", command)
		}
	}

	// Manejo de errores.
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error leyendo la entrada:", err)
	}
}
