package analyzer

import (
	"bufio"
	"flag"
	"fmt"
	"godisk/commands"
	"io"
	"os"
	"strings"
)

// ProcessCommands recibe un string con comandos y los procesa línea por línea
// Devuelve la salida completa como string
func ProcessCommands(input string) string {
	var outputBuilder strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(input))

	for scanner.Scan() {
		line := scanner.Text()

		// Ignora líneas vacías
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Si es un comentario, ignorarlo
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		// Procesa la línea de comando actual
		outputBuilder.WriteString(fmt.Sprintf("> %s\n", line))

		// Si el usuario quiere salir, retornamos inmediatamente
		if strings.ToLower(line) == "exit" {
			outputBuilder.WriteString("Saliendo...\n")
			continue
		}

		// Ejecuta el comando y captura su salida
		output := executeCommand(line)
		outputBuilder.WriteString(output)
		outputBuilder.WriteString("\n")
	}

	return outputBuilder.String()
}

// executeCommand ejecuta un comando específico con sus argumentos
func executeCommand(commandLine string) string {
	// Divide la línea en partes (comando y argumentos)
	parts := strings.Fields(commandLine)
	if len(parts) == 0 {
		return ""
	}

	// La primera palabra es el comando
	command := strings.ToLower(parts[0])
	args := parts[1:]

	// Guarda stdout original para restaurarlo después
	oldStdout := os.Stdout

	// Crea un pipe para capturar stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Ejecuta el comando según su tipo
	switch command {
	case "mkdisk":
		mkdiskCmd := flag.NewFlagSet("mkdisk", flag.ContinueOnError)
		size := mkdiskCmd.Int("size", 0, "Tamaño del disco.")
		unit := mkdiskCmd.String("unit", "m", "Unidad del tamaño (k/m).")
		fit := mkdiskCmd.String("fit", "ff", "Tipo de ajuste (bf/ff/wf).")
		path := mkdiskCmd.String("path", "", "Ruta del disco a crear.")

		mkdiskCmd.Parse(args)
		commands.ExecuteMkdisk(*size, *unit, *fit, *path)

	case "rmdisk":
		rmdiskCmd := flag.NewFlagSet("rmdisk", flag.ContinueOnError)
		path := rmdiskCmd.String("path", "", "Ruta del disco a eliminar.")
		rmdiskCmd.Parse(args)
		commands.ExecuteRmdisk(*path)

	case "fdisk":
		fdiskCmd := flag.NewFlagSet("fdisk", flag.ContinueOnError)
		size := fdiskCmd.Int64("size", 0, "Tamaño de la partición.")
		path := fdiskCmd.String("path", "", "Ruta del disco.")
		name := fdiskCmd.String("name", "", "Nombre de la partición.")
		unit := fdiskCmd.String("unit", "k", "Unidad del tamaño (b/k/m).")
		typeStr := fdiskCmd.String("type", "p", "Tipo de partición (p/e/l).")
		fit := fdiskCmd.String("fit", "wf", "Tipo de ajuste (bf/ff/wf).")

		fdiskCmd.Parse(args)
		commands.ExecuteFdisk(*path, *name, *unit, *typeStr, *fit, *size)

	case "mount":
		mountCmd := flag.NewFlagSet("mount", flag.ContinueOnError)
		path := mountCmd.String("path", "", "Ruta del disco.")
		name := mountCmd.String("name", "", "Nombre de la partición.")
		mountCmd.Parse(args)
		commands.ExecuteMount(*path, *name)

	case "mounted":
		commands.ExecuteMounted()

	case "mkfs":
		mkfsCmd := flag.NewFlagSet("mkfs", flag.ContinueOnError)
		id := mkfsCmd.String("id", "", "ID de la partición a formatear.")
		typeStr := mkfsCmd.String("type", "full", "Tipo de formateo (full).")
		fs := mkfsCmd.String("fs", "2fs", "Sistema de archivos (2fs).")

		mkfsCmd.Parse(args)
		commands.ExecuteMkfs(*id, *typeStr, *fs)

	case "rep":
		repCmd := flag.NewFlagSet("rep", flag.ContinueOnError)
		name := repCmd.String("name", "", "Nombre del reporte (mbr, disk).")
		path := repCmd.String("path", "", "Ruta donde se guardará el reporte.")
		id := repCmd.String("id", "", "ID de la partición montada.")
		repCmd.Parse(args)
		commands.ExecuteRep(*name, *path, *id)

	case "login":
		loginCmd := flag.NewFlagSet("login", flag.ContinueOnError)
		user := loginCmd.String("user", "", "Nombre de usuario")
		pass := loginCmd.String("pass", "", "Contraseña")
		id := loginCmd.String("id", "", "ID de la partición montada")

		loginCmd.Parse(args)
		commands.ExecuteLogin(*user, *pass, *id)

	case "logout":
		commands.ExecuteLogout()
	case "mkfile":
		mkfileCmd := flag.NewFlagSet("mkfile", flag.ContinueOnError)
		path := mkfileCmd.String("path", "", "Ruta del archivo a crear.")
		rFlag := mkfileCmd.Bool("r", false, "Crear carpetas padres si no existen.")
		size := mkfileCmd.Int("size", 0, "Tamaño en bytes del archivo.")
		cont := mkfileCmd.String("cont", "", "Ruta del archivo con contenido.")
		// Ignorar errores de parsing para que no termine el programa
		_ = mkfileCmd.Parse(args)
		commands.ExecuteMkfile(*path, *rFlag, *size, *cont)

	default:
		return fmt.Sprintf("Comando '%s' no reconocido.\n", command)
	}

	// Cierra el pipe de escritura para poder leer todo
	w.Close()

	// Lee la salida capturada
	var buf strings.Builder
	io.Copy(&buf, r)

	// Restaura stdout
	os.Stdout = oldStdout

	return buf.String()
}
