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
	// Crea un nuevo scanner que lee desde la entrada estándar (teclado)
	// bufio.Scanner lee línea por línea de manera eficiente y maneja automáticamente los buffers
	scanner := bufio.NewScanner(os.Stdin)

	// Bucle infinito para mantener el programa en ejecución hasta que el usuario decida salir
	// Este patrón se conoce como REPL (Read-Eval-Print Loop)
	for {
		// Muestra el prompt ">" para indicar al usuario que puede introducir un comando
		fmt.Print("> ")

		// Intenta leer una línea completa de entrada del usuario
		// scanner.Scan() retorna true si logró leer una línea, false si hubo error o EOF
		if !scanner.Scan() {
			break // Si no puede leer (EOF o error), sale del bucle
		}

		// Obtiene el texto de la línea leída, sin el carácter de nueva línea
		line := scanner.Text()

		// Verifica si el usuario quiere salir del programa
		// strings.ToLower() convierte a minúsculas para hacer la comparación insensible a mayúsculas
		if strings.ToLower(line) == "exit" {
			fmt.Println("Saliendo...")
			break // Termina el bucle y por tanto el programa
		}

		// Divide la línea en palabras separadas por espacios en blanco
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
			// Crear un nuevo conjunto de flags específico para el comando mkdisk
			// flag.ExitOnError hace que el programa termine si hay un error en los flags
			mkdiskCmd := flag.NewFlagSet("mkdisk", flag.ExitOnError)

			// Definir los flags que acepta el comando mkdisk:
			// Int() crea un flag que acepta números enteros
			size := mkdiskCmd.Int("size", 0, "Tamaño del disco.")

			// String() crea un flag que acepta cadenas de texto
			unit := mkdiskCmd.String("unit", "m", "Unidad del tamaño (k/m).")
			fit := mkdiskCmd.String("fit", "ff", "Tipo de ajuste (bf/ff/wf).")
			path := mkdiskCmd.String("path", "", "Ruta del disco a crear.")

			// Procesa los argumentos pasados al comando y asigna valores a los flags
			// Parse() analiza los argumentos y llena las variables con los valores correspondientes
			mkdiskCmd.Parse(args)

			// Validaciones obligatorias antes de ejecutar el comando

			// El parámetro path es obligatorio, no puede estar vacío
			// *path desreferencia el puntero para obtener el valor real
			if *path == "" {
				fmt.Println("Error: el parámetro -path es obligatorio para mkdisk.")
				continue // Vuelve al inicio del bucle sin ejecutar el comando
			}

			// El parámetro size debe ser positivo y mayor que cero
			if *size <= 0 {
				fmt.Println("Error: el parámetro -size es obligatorio y debe ser positivo.")
				continue // Vuelve al inicio del bucle sin ejecutar el comando
			}

			// Si todas las validaciones pasan, ejecuta el comando mkdisk
			// Pasa los valores desreferenciados (con *) a la función
			commands.ExecuteMkdisk(*size, *unit, *fit, *path)

		default:
			// Si el comando no es reconocido, mostrar mensaje de error
			fmt.Printf("Comando '%s' no reconocido.\n", command)
		}
	}

	// Manejo de errores del scanner
	// scanner.Err() retorna cualquier error que haya ocurrido durante la lectura
	if err := scanner.Err(); err != nil {
		// fmt.Fprintln() imprime en stderr (salida de error estándar) en lugar de stdout
		fmt.Fprintln(os.Stderr, "Error leyendo la entrada:", err)
	}
}
