package main

import "fmt"

// Estructura para representar un Producto
type Producto struct {
	ID       int
	Nombre   string
	Cantidad int
	Precio   float64
}

// Método para mostrar información del producto
func (p Producto) MostrarDetalles() {
	fmt.Printf("ID: %d\nNombre: %s\nCantidad: %d\nPrecio: %.2f\n", p.ID, p.Nombre, p.Cantidad, p.Precio)
}

// Método para actualizar la cantidad del producto
func (p *Producto) ActualizarCantidad(nuevaCantidad int) {
	p.Cantidad = nuevaCantidad
}

// Método para calcular el valor total del inventario del producto
func (p Producto) CalcularValorTotal() float64 {
	return float64(p.Cantidad) * p.Precio
}

// Función principal
func main3() {
	// Crear una lista de productos
	inventario := []Producto{
		{ID: 1, Nombre: "Laptop", Cantidad: 10, Precio: 799.99},
		{ID: 2, Nombre: "Teléfono", Cantidad: 20, Precio: 599.50},
		{ID: 3, Nombre: "Tablet", Cantidad: 15, Precio: 399.00},
	}

	// Mostrar detalles de cada producto
	fmt.Println("Inventario inicial:")
	for _, producto := range inventario {
		producto.MostrarDetalles()
		fmt.Printf("Valor total: %.2f\n\n", producto.CalcularValorTotal())
	}

	// Actualizar la cantidad de un producto
	fmt.Println("Actualizando el inventario...")
	inventario[0].ActualizarCantidad(8)

	// Mostrar el inventario actualizado
	fmt.Println("\nInventario actualizado:")
	for _, producto := range inventario {
		producto.MostrarDetalles()
		fmt.Printf("Valor total: %.2f\n\n", producto.CalcularValorTotal())
	}
}
