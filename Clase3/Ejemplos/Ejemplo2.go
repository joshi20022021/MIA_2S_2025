package main

import (
	"errors"
	"fmt"
	"math"
)

// Función para calcular el interés compuesto
func calcularInteresCompuesto(principal float64, tasa float64, tiempo int) (float64, error) {
	if principal < 0 || tasa < 0 || tiempo < 0 {
		return 0, errors.New("los valores no pueden ser negativos")
	}
	monto := principal * math.Pow(1+(tasa/100), float64(tiempo))
	return monto, nil
}

// Función principal
func main2() {
	// Solicitar datos al usuario
	var principal float64
	var tasa float64
	var tiempo int

	fmt.Println("Calculadora de Interés Compuesto")
	fmt.Print("Introduce el monto principal: ")
	fmt.Scan(&principal)
	fmt.Print("Introduce la tasa de interés (en porcentaje): ")
	fmt.Scan(&tasa)
	fmt.Print("Introduce el tiempo (en años): ")
	fmt.Scan(&tiempo)

	// Calcular y manejar posibles errores
	montoFinal, err := calcularInteresCompuesto(principal, tasa, tiempo)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Mostrar resultados
	fmt.Printf("\nResultados:\n")
	fmt.Printf("Monto principal: %.2f\n", principal)
	fmt.Printf("Tasa de interés: %.2f%%\n", tasa)
	fmt.Printf("Tiempo: %d años\n", tiempo)
	fmt.Printf("Monto final: %.2f\n", montoFinal)
}
