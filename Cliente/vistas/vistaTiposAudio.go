/**
 * @file vistaTiposAudio.go
 * @brief Vista que muestra los tipos de audio disponibles.
 */
package vistas

import "cliente.local/cliente/utilidades"

/**
 * @brief Muestra los tipos de audio y permite seleccionar uno.
 *
 * Retorna al menú principal cuando el usuario elige "Atrás", cuando abandona
 * una reproducción o si no es posible consultar los tipos.
 */
func (thisN *NavegadorVistas) mostrarTiposAudio() {
	for {
		tipos, err := thisN.controladorAudios.ObtenerTiposAudio()
		if err != nil {
			thisN.avisar("No se pudieron consultar los tipos de audio: %v", err)
			return
		}

		thisN.imprimirEncabezado()
		nombres := make([]string, len(tipos))
		for indice, tipo := range tipos {
			nombres[indice] = tipo.Nombre
		}
		opcionAtras := imprimirOpcionesConAtras(nombres)

		opcion := utilidades.LeerOpcion("Seleccione una opción: ", 1, opcionAtras)
		if opcion == opcionAtras {
			return
		}

		if volverAlMenuPrincipal := thisN.mostrarListaAudios(tipos[opcion-1]); volverAlMenuPrincipal {
			return
		}
	}
}
