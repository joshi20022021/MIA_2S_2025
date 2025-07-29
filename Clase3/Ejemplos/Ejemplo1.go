package main

import (
	"fmt"
	"time"
)

// Función principal
func main() {
	// Declaracion variables
	var nombre string = "Carlos"
	edad := 30
	ciudad := "Madrid"

	//mensaje de bienvenida
	fmt.Printf("Bienvenido, %s. Tienes %d años y vives en %s.\n", nombre, edad, ciudad)

	// Operaciones matemáticas
	a, b := 12, 8
	fmt.Printf("Operaciones básicas con %d y %d:\n", a, b)
	fmt.Printf("Suma: %d\n", a+b)
	fmt.Printf("Resta: %d\n", a-b)
	fmt.Printf("Multiplicación: %d\n", a*b)
	fmt.Printf("División: %.2f\n", float64(a)/float64(b))

	// Condicional y formato de hora
	horaActual := time.Now()
	fmt.Println("\nEstado del día:")
	if horaActual.Hour() < 12 {
		fmt.Println("Es la mañana. ¡Buenos días!")
	} else if horaActual.Hour() < 18 {
		fmt.Println("Es la tarde. ¡Sigue adelante!")
	} else {
		fmt.Println("Es la noche. ¡Hora de descansar!")
	}

	// Múltiples variables con valores intercambiados
	x, y := 10, 20
	fmt.Printf("\nAntes del intercambio: x = %d, y = %d\n", x, y)
	x, y = y, x
	fmt.Printf("Después del intercambio: x = %d, y = %d\n", x, y)
}
