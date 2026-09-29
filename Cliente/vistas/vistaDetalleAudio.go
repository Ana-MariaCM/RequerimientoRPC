/**
 * @file vistaDetalleAudio.go
 * @brief Vista que muestra los detalles (metadatos) de un audio.
 */
package vistas

import (
	"fmt"

	dtos "cliente.local/cliente/capaFachadaServices/DTOs"
	"cliente.local/cliente/utilidades"
)

/** @brief Opción de la vista de detalle para reproducir el audio. */
const opcionReproducir = 1

/** @brief Opción de la vista de detalle para volver a la lista de audios. */
const opcionAtrasDetalle = 2

/**
 * @brief Muestra los metadatos de un audio y permite reproducirlo.
 * @param idAudio Identificador del audio seleccionado.
 * @return true si se reprodujo el audio (se vuelve al menú principal),
 *         false si el usuario eligió "Atrás".
 */
func (thisN *NavegadorVistas) mostrarDetalleAudio(idAudio int) bool {
	detalle, err := thisN.controladorAudios.ObtenerDetalleAudio(idAudio)
	if err != nil {
		thisN.avisar("No se pudo consultar el detalle del audio: %v", err)
		return false
	}

	for {
		thisN.imprimirEncabezado()
		imprimirTituloAudio(detalle)
		for _, metadato := range detalle.Metadatos {
			fmt.Printf("  • %s: %s\n", metadato.Nombre, metadato.Valor)
		}
		fmt.Println()
		fmt.Println("  1. Reproducir audio")
		fmt.Println("  2. Atrás")
		fmt.Println()

		if utilidades.LeerOpcion("Seleccione una opción: ", opcionReproducir, opcionAtrasDetalle) == opcionAtrasDetalle {
			return false
		}

		if thisN.mostrarReproduccion(detalle) {
			return true
		}
	}
}

/**
 * @brief Imprime la línea "Tipo: Título - archivo" que identifica al audio.
 * @param detalle Detalle del audio.
 */
func imprimirTituloAudio(detalle dtos.AudioDetalleDTO) {
	fmt.Printf("%s: %s - %s\n\n", detalle.NombreTipo, detalle.Titulo, detalle.NombreArchivo)
}
