//•	Imprimir un encabezado: ==== CALCULADORA CIENTÍFICA v1.0 ====
// •	Declarar dos variables 'a' y 'b' con valores fijos (por ejemplo, 10 y 3).
// •	Imprimir el resultado de la suma en formato legible: a + b = 13

package main

import (
	"fmt"
)

func main() {

	var historial string    // Variable para almacenar el historial
	var contador int        // Variable para contar el total de operaciones
	var seguir string = "s" // Variable para controlar el ciclo

	for seguir == "s" { // Ciclo para repetir el programa
		var a, b float64
		var operacion string

		fmt.Println("\n==== CALCULADORA ====")
		fmt.Println("Operaciones: +, -, *, /, ^, !")
		fmt.Print("Ingrese operación: ")
		fmt.Scan(&operacion)

		// FACTORIAL
		if operacion == "!" {
			var n int
			fmt.Print("Ingrese un número: ")
			fmt.Scan(&n)

			resultado := 1
			for i := 1; i <= n; i++ {
				resultado *= i
			}
			fmt.Println("Resultado:", resultado)
			contador++
			historial += fmt.Sprintf("%d) %d ! = %d\n", contador, n, resultado)

		} else {
			fmt.Print("Ingrese el primer número: ")
			fmt.Scan(&a)

			fmt.Print("Ingrese el segundo número: ")
			fmt.Scan(&b)

			var resultado float64

			switch operacion {
			case "+":
				resultado = a + b
			case "-":
				resultado = a - b
			case "*":
				resultado = a * b

			case "/":
				if b != 0 {
					resultado = a / b
				} else {
					fmt.Println("Error: división por 0")
					continue
				}
			//potencia
			case "^":
				resultado = a
				for i := 1; i < int(b); i++ {
					resultado *= a
				}

			default:
				fmt.Println("Operación inválida")
				continue
			}

			fmt.Println("Resultado:", resultado)

			contador++
			historial += fmt.Sprintf("%d) %.2f %s %.2f = %.2f\n", contador, a, operacion, b, resultado)
		}

		fmt.Print("¿Otra operación? (s/n): ")
		fmt.Scan(&seguir)
	}

	//  HISTORIAL
	fmt.Println("\n==== HISTORIAL DE OPERACIONES ====")
	fmt.Print(historial)
	fmt.Println("Total de operaciones:", contador)
	fmt.Println("¡Hasta luego!")
}
