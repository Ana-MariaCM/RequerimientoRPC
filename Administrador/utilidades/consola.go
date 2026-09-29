/**
 * @file consola.go
 * @brief Funciones de apoyo para la lectura de datos por consola.
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
 * @brief Lee un texto que no puede quedar vacío.
 * @param mensaje Texto que se muestra antes de leer.
 * @return Texto ingresado.
 */
func LeerTextoObligatorio(mensaje string) string {
	for {
		texto := LeerTexto(mensaje)
		if texto != "" {
			return texto
		}
		fmt.Println("Este dato es obligatorio.")
	}
}

/**
 * @brief Lee un número entero obligatorio.
 * @param mensaje Texto que se muestra antes de leer.
 * @return El número ingresado, en formato texto.
 */
func LeerNumeroObligatorio(mensaje string) string {
	for {
		texto := LeerTextoObligatorio(mensaje)
		if _, err := strconv.Atoi(texto); err == nil {
			return texto
		}
		fmt.Println("Ingrese un número entero.")
	}
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
 * @brief Convierte una cantidad de bytes en un texto legible (B, KB o MB).
 * @param bytes Cantidad de bytes.
 * @return Texto con la unidad adecuada.
 */
func FormatearTamanio(bytes int64) string {
	switch {
	case bytes >= 1<<20:
		return fmt.Sprintf("%.2f MB", float64(bytes)/(1<<20))
	case bytes >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(bytes)/(1<<10))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
