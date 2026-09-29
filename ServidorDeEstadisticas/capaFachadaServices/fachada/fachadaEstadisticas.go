/**
 * @file fachadaEstadisticas.go
 * @brief Fachada que registra las reproducciones y calcula las estadísticas.
 */
package fachada

import (
	capaaccesodatos "estadisticas/capaAccesoADatos"
	dtos "estadisticas/capaFachadaServices/DTOs"
	"estadisticas/modelos"
)

/**
 * @brief Fachada del servidor de estadísticas.
 */
type FachadaEstadisticas struct {
	repositorio *capaaccesodatos.RepositorioEstadisticas ///< Acceso a las reproducciones almacenadas.
}

/**
 * @brief Crea la fachada de estadísticas.
 * @return Puntero a la nueva fachada.
 */
func NuevaFachadaEstadisticas() *FachadaEstadisticas {
	return &FachadaEstadisticas{repositorio: capaaccesodatos.GetRepositorioEstadisticas()}
}

/**
 * @brief Registra una reproducción recibida desde la cola.
 * @param reproduccionDTO Datos de la reproducción deserializados del mensaje JSON.
 * @return La reproducción almacenada y el resumen actualizado de estadísticas.
 */
func (thisF *FachadaEstadisticas) RegistrarReproduccion(reproduccionDTO dtos.ReproduccionDTOInput) (modelos.Reproduccion, modelos.ResumenEstadisticas) {
	reproduccion := modelos.Reproduccion{
		IdAudio:          reproduccionDTO.IdAudio,
		Titulo:           reproduccionDTO.Titulo,
		Tipo:             reproduccionDTO.Tipo,
		NombreArchivo:    reproduccionDTO.NombreArchivo,
		Usuario:          reproduccionDTO.Usuario,
		DireccionCliente: reproduccionDTO.DireccionCliente,
		FechaHora:        reproduccionDTO.FechaHora,
	}

	almacenada := thisF.repositorio.GuardarReproduccion(reproduccion)
	return almacenada, thisF.repositorio.CalcularResumen()
}
