/**
 * @file vistaListaAudios.go
 * @brief Vista que muestra la lista de audios de un tipo.
 */
package vistas

import (
	"fmt"

	dtos "cliente.local/cliente/capaFachadaServices/DTOs"
	"cliente.local/cliente/utilidades"
)

/**
 * @brief Muestra los audios de un tipo y permite seleccionar uno.
 * @param tipo Tipo de audio seleccionado.
 * @return true si se debe volver al menú principal (tras una reproducción),
 *         false si el usuario eligió "Atrás".
 */
func (thisN *NavegadorVistas) mostrarListaAudios(tipo dtos.TipoAudioDTO) bool {
	for {
		audios, err := thisN.controladorAudios.ObtenerAudiosPorTipo(tipo.Id)
		if err != nil {
			thisN.avisar("No se pudieron consultar los audios de %s: %v", tipo.Nombre, err)
			return false
		}

		thisN.imprimirEncabezado()
		fmt.Printf("Tipo: %s\n\n", tipo.Nombre)
		titulos := make([]string, len(audios))
		for indice, audio := range audios {
			titulos[indice] = audio.Titulo
		}
		opcionAtras := imprimirOpcionesConAtras(titulos)

		opcion := utilidades.LeerOpcion("Seleccione una opción: ", 1, opcionAtras)
		if opcion == opcionAtras {
			return false
		}

		if volverAlMenuPrincipal := thisN.mostrarDetalleAudio(audios[opcion-1].Id); volverAlMenuPrincipal {
			return true
		}
	}
}
