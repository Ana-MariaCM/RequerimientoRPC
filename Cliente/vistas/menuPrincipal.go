/**
 * @file menuPrincipal.go
 * @brief Vista del menú principal del cliente.
 */
package vistas

import (
	"fmt"

	"cliente.local/cliente/utilidades"
)

/** @brief Opción del menú principal para ver los tipos de audio. */
const opcionVerTiposAudio = 1

/** @brief Opción del menú principal para salir de la aplicación. */
const opcionSalir = 2

/**
 * @brief Muestra el menú principal hasta que el usuario elija salir.
 *
 * Menú:
 *  1. Ver tipos de audio
 *  2. Salir
 */
func (thisN *NavegadorVistas) MostrarMenuPrincipal() {
	for {
		thisN.imprimirEncabezado()
		fmt.Println("  1. Ver tipos de audio")
		fmt.Println("  2. Salir")
		fmt.Println()

		switch utilidades.LeerOpcion("Seleccione una opción: ", opcionVerTiposAudio, opcionSalir) {
		case opcionVerTiposAudio:
			thisN.mostrarTiposAudio()
		case opcionSalir:
			fmt.Println("Hasta pronto.")
			return
		}
	}
}
