/**
 * @file vistaEstadisticas.go
 * @brief Vista que imprime por pantalla las reproducciones y las estadísticas.
 */
package vistas

import (
	"fmt"
	"strings"

	"estadisticas/modelos"
)

/**
 * @brief Imprime los datos de una reproducción recién registrada.
 * @param reproduccion Reproducción a mostrar.
 */
func MostrarReproduccion(reproduccion modelos.Reproduccion) {
	fmt.Printf("  Reproducción #%d\n", reproduccion.Numero)
	fmt.Printf("    Audio        : %s (id %d)\n", reproduccion.Titulo, reproduccion.IdAudio)
	fmt.Printf("    Tipo         : %s\n", reproduccion.Tipo)
	fmt.Printf("    Archivo      : %s\n", reproduccion.NombreArchivo)
	fmt.Printf("    Usuario      : %s (%s)\n", reproduccion.Usuario, reproduccion.DireccionCliente)
	fmt.Printf("    Fecha y hora : %s\n", reproduccion.FechaHora)
}

/**
 * @brief Imprime el resumen de estadísticas acumuladas.
 * @param resumen Resumen calculado por la fachada.
 */
func MostrarResumen(resumen modelos.ResumenEstadisticas) {
	fmt.Println("  " + strings.Repeat("-", 50))
	fmt.Printf("  Estadísticas acumuladas: %d reproducciones\n", resumen.TotalReproducciones)
	mostrarConteos("Por tipo de audio", resumen.PorTipo)
	mostrarConteos("Por audio", resumen.PorAudio)
	mostrarConteos("Por usuario", resumen.PorUsuario)
	fmt.Println("  " + strings.Repeat("-", 50))
}

/**
 * @brief Imprime una tabla de conteos con un título.
 * @param titulo Título de la tabla.
 * @param conteos Conteos a imprimir.
 */
func mostrarConteos(titulo string, conteos []modelos.Conteo) {
	fmt.Printf("  %s:\n", titulo)
	for _, conteo := range conteos {
		fmt.Printf("    %-50s %3d\n", conteo.Nombre, conteo.Cantidad)
	}
}
