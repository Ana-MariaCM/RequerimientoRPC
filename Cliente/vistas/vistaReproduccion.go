/**
 * @file vistaReproduccion.go
 * @brief Vista que se muestra mientras un audio se reproduce por streaming.
 */
package vistas

import (
	"fmt"
	"strings"

	dtos "cliente.local/cliente/capaFachadaServices/DTOs"
	"cliente.local/cliente/utilidades"
)

/** @brief Opción de la vista de reproducción para abandonar el audio. */
const opcionSalirReproduccion = 1

/**
 * @brief Inicia la reproducción y espera a que el usuario la abandone.
 *
 * El usuario puede elegir "Salir" en cualquier momento, incluso antes de que
 * termine el audio; en ese caso se detiene el streaming y se vuelve al menú
 * principal.
 * @param detalle Audio a reproducir.
 * @return true si la reproducción inició (se debe volver al menú principal),
 *         false si no pudo iniciarse (se vuelve al detalle del audio).
 */
func (thisN *NavegadorVistas) mostrarReproduccion(detalle dtos.AudioDetalleDTO) bool {
	sesion, err := thisN.controladorReproduccion.IniciarReproduccion(detalle)
	if err != nil {
		thisN.avisar("No se pudo reproducir el audio: %v", err)
		return false
	}

	thisN.imprimirEncabezado()
	imprimirTituloAudio(detalle)
	fmt.Printf("%sReproduciendo audio\n\n", strings.Repeat(" ", 12))
	fmt.Println("  1. Salir")
	fmt.Println()

	usuarioSalio := make(chan struct{})
	go notificarFinDeReproduccion(sesion.Finalizada(), sesion.Error, usuarioSalio)

	utilidades.LeerOpcion("Seleccione una opción: ", opcionSalirReproduccion, opcionSalirReproduccion)
	close(usuarioSalio)
	sesion.Detener()
	return true
}

/**
 * @brief Informa al usuario cuando el audio termina antes de que él elija "Salir".
 * @param finalizada Canal que se cierra al terminar la reproducción.
 * @param obtenerError Función que devuelve el error de la reproducción (si lo hubo).
 * @param usuarioSalio Canal que se cierra cuando el usuario abandona la reproducción.
 */
func notificarFinDeReproduccion(finalizada <-chan struct{}, obtenerError func() error, usuarioSalio <-chan struct{}) {
	select {
	case <-usuarioSalio:
	case <-finalizada:
		select {
		case <-usuarioSalio:
			return
		default:
		}
		if err := obtenerError(); err != nil {
			fmt.Printf("\n  La reproducción se interrumpió: %v\n", err)
		} else {
			fmt.Println("\n  Reproducción finalizada.")
		}
		fmt.Print("Seleccione 1 para volver al menú principal: ")
	}
}
