/**
 * @file navegadorVistas.go
 * @brief Estructura común de las vistas del cliente y funciones de presentación compartidas.
 */
package vistas

import (
	"fmt"
	"strings"

	capacontroladores "cliente.local/cliente/capaControladores"
	"cliente.local/cliente/utilidades"
)

/** @brief Ancho (en caracteres) de los encabezados de las pantallas. */
const anchoPantalla = 60

/**
 * @brief Agrupa los controladores que usan las vistas y el aviso pendiente de mostrar.
 */
type NavegadorVistas struct {
	controladorAudios       *capacontroladores.ControladorAudios       ///< Consultas de metadatos.
	controladorReproduccion *capacontroladores.ControladorReproduccion ///< Reproducción por streaming.
	aviso                   string                                     ///< Mensaje que se mostrará en la siguiente pantalla.
}

/**
 * @brief Crea el navegador de vistas del cliente.
 * @param controladorAudios Controlador de consultas de metadatos.
 * @param controladorReproduccion Controlador de reproducción.
 * @return Puntero al navegador.
 */
func NuevoNavegadorVistas(controladorAudios *capacontroladores.ControladorAudios,
	controladorReproduccion *capacontroladores.ControladorReproduccion) *NavegadorVistas {
	return &NavegadorVistas{controladorAudios: controladorAudios, controladorReproduccion: controladorReproduccion}
}

/**
 * @brief Limpia la pantalla, imprime el encabezado "Spotify" y el aviso pendiente.
 */
func (thisN *NavegadorVistas) imprimirEncabezado() {
	utilidades.LimpiarPantalla()
	linea := strings.Repeat("=", anchoPantalla)
	titulo := "Spotify"
	fmt.Println(linea)
	fmt.Printf("%s%s\n", strings.Repeat(" ", (anchoPantalla-len(titulo))/2), titulo)
	fmt.Println(linea)
	if thisN.aviso != "" {
		fmt.Printf("(!) %s\n\n", thisN.aviso)
		thisN.aviso = ""
	}
}

/**
 * @brief Registra un mensaje para mostrarlo en la siguiente pantalla.
 * @param formato Formato estilo printf.
 * @param argumentos Argumentos del formato.
 */
func (thisN *NavegadorVistas) avisar(formato string, argumentos ...any) {
	thisN.aviso = fmt.Sprintf(formato, argumentos...)
}

/**
 * @brief Imprime una lista numerada de opciones seguida de la opción "Atrás".
 * @param opciones Textos de las opciones.
 * @return Número de la opción "Atrás".
 */
func imprimirOpcionesConAtras(opciones []string) int {
	for indice, opcion := range opciones {
		fmt.Printf("  %d. %s\n", indice+1, opcion)
	}
	opcionAtras := len(opciones) + 1
	fmt.Printf("  %d. Atrás\n\n", opcionAtras)
	return opcionAtras
}
