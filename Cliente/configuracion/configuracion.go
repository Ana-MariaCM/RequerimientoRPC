/**
 * @file configuracion.go
 * @brief Parámetros de configuración del cliente.
 */
package configuracion

import (
	"os"
	"strings"
)

/** @brief URL por defecto del servidor de metadatos (REST). */
const URLMetadatosPorDefecto = "http://localhost:5001"

/** @brief Dirección por defecto del servidor de streaming (gRPC). */
const DireccionStreamingPorDefecto = "localhost:50051"

/**
 * @brief Obtiene la URL base del servidor de metadatos (variable URL_METADATOS).
 * @return URL base sin barra final.
 */
func ObtenerURLMetadatos() string {
	return strings.TrimRight(obtenerVariable("URL_METADATOS", URLMetadatosPorDefecto), "/")
}

/**
 * @brief Obtiene la dirección host:puerto del servidor de streaming (variable DIRECCION_STREAMING).
 * @return Dirección del servidor gRPC.
 */
func ObtenerDireccionStreaming() string {
	return obtenerVariable("DIRECCION_STREAMING", DireccionStreamingPorDefecto)
}

/**
 * @brief Lee una variable de entorno y devuelve un valor por defecto si está vacía.
 * @param nombre Nombre de la variable de entorno.
 * @param valorPorDefecto Valor que se retorna si la variable no está definida.
 * @return Valor de la variable o el valor por defecto.
 */
func obtenerVariable(nombre string, valorPorDefecto string) string {
	valor := os.Getenv(nombre)
	if valor == "" {
		return valorPorDefecto
	}
	return valor
}
