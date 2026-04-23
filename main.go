package main

import (
	"errors"
	"fmt"
	"semana03-taller-relaciones/internal/cafeteria"
)

func main() {
	var repo cafeteria.Repository = cafeteria.NewRepoMemoria()

	catBebidas := cafeteria.Categoria{ID: 1, Nombre: "Bebidas"}
	catSnacks := cafeteria.Categoria{ID: 2, Nombre: "Snacks"}

	_ = repo.GuardarCliente(cafeteria.Cliente{ID: 1, Nombre: "Ana", Carrera: "TI", Saldo: 50})
	_ = repo.GuardarCliente(cafeteria.Cliente{ID: 2, Nombre: "Luis", Carrera: "Sistemas", Saldo: 30})

	_ = repo.GuardarProducto(cafeteria.Producto{ID: 1, Nombre: "Café", Precio: 1.25, Stock: 40, Categoria: catBebidas})
	_ = repo.GuardarProducto(cafeteria.Producto{ID: 2, Nombre: "Jugo", Precio: 2.00, Stock: 25, Categoria: catBebidas})
	_ = repo.GuardarProducto(cafeteria.Producto{ID: 3, Nombre: "Sandwich", Precio: 4.50, Stock: 15, Categoria: catSnacks})

	// Cliente que existe
	c, err := repo.ObtenerCliente(1)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Cliente encontrado:", c.Nombre)
	}

	// Cliente que NO existe
	c, err = repo.ObtenerCliente(99)
	if err != nil {
		fmt.Println("Error:", err)
		if errors.Is(err, cafeteria.ErrClienteNoEncontrado) {
			fmt.Println("(Error confirmado: ErrClienteNoEncontrado)")
		}
	} else {
		fmt.Println("Cliente encontrado:", c.Nombre)
	}

	fmt.Println("\n--- Listado de productos ---")
	for _, p := range repo.ListarProductos() {
		fmt.Printf("- %s | $%.2f | stock %d | categoría: %s\n",
			p.Nombre, p.Precio, p.Stock, p.Categoria.Nombre)
	}

	// Pedido con Cliente y Producto completos (obtenidos del repo, no solo IDs)
	clientePedido, _ := repo.ObtenerCliente(2)
	productoPedido, _ := repo.ObtenerProducto(3)
	pedido := cafeteria.Pedido{
		ID:       1,
		Cliente:  clientePedido,
		Producto: productoPedido,
		Cantidad: 1,
		Total:    productoPedido.Precio * 1,
		Fecha:    "2026-04-23",
	}

	fmt.Println("\n--- Pedido (structs completos anidados) ---")
	fmt.Printf("%+v\n", pedido)
	fmt.Printf("  → Cliente completo: %+v\n", pedido.Cliente)
	fmt.Printf("  → Producto completo (incl. categoría): %+v\n", pedido.Producto)
}
