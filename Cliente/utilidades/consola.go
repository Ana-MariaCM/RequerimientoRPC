/**
 * @file consola.go
 * @brief Funciones de apoyo para la interacción por consola.
 */
package utilidades

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/** @brief Lector compartido de la entrada estándar. */
var lectorEntrada = bufio.NewReader(os.Stdin)

/**
 * @brief Muestra un mensaje y lee una línea de texto.
 * @param mensaje Texto que se muestra antes de leer.
 * @return Texto ingresado sin espacios al inicio ni al final.
 */
func LeerTexto(mensaje string) string {
	fmt.Print(mensaje)
	linea, err := lectorEntrada.ReadString('\n')
	if err != nil && linea == "" {
		// Fin de la entrada estándar: se termina la aplicación.
		fmt.Println()
		os.Exit(0)
	}
	return strings.TrimSpace(linea)
}

/**
 * @brief Lee una opción numérica comprendida entre un mínimo y un máximo.
 *
 * Repite la lectura hasta que el usuario ingrese un número válido.
 * @param mensaje Texto que se muestra antes de leer.
 * @param minimo Valor mínimo permitido.
 * @param maximo Valor máximo permitido.
 * @return Opción seleccionada.
 */
func LeerOpcion(mensaje string, minimo int, maximo int) int {
	for {
		opcion, err := strconv.Atoi(LeerTexto(mensaje))
		if err == nil && opcion >= minimo && opcion <= maximo {
			return opcion
		}
		fmt.Printf("Opción no válida, ingrese un número entre %d y %d.\n", minimo, maximo)
	}
}

/**
 * @brief Limpia la pantalla de la terminal (secuencia ANSI).
 */
func LimpiarPantalla() {
	fmt.Print("\033[H\033[2J")
}
