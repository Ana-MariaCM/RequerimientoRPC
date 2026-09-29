/**
 * @file configuracion.go
 * @brief Parámetros de configuración del servidor de audios.
 */
package configuracion

import (
	"os"
	"strings"
)

/** @brief Puerto HTTP por defecto del servidor REST de audios. */
const PuertoPorDefecto = "5000"

/** @brief Carpeta por defecto donde se almacenan los archivos mp3. */
const RutaAudiosPorDefecto = "audios"

/** @brief URL por defecto del servidor de metadatos donde se registran los audios nuevos. */
const URLMetadatosPorDefecto = "http://localhost:5001"

/** @brief Tamaño máximo en bytes de un audio recibido (100 MB). */
const TamanioMaximoAudio = 100 << 20

/**
 * @brief Obtiene el puerto del servidor REST (variable PUERTO_AUDIOS).
 * @return Puerto en formato texto.
 */
func ObtenerPuerto() string {
	return obtenerVariable("PUERTO_AUDIOS", PuertoPorDefecto)
}

/**
 * @brief Obtiene la carpeta donde se almacenan los audios (variable RUTA_AUDIOS).
 * @return Ruta de la carpeta de audios.
 */
func ObtenerRutaAudios() string {
	return obtenerVariable("RUTA_AUDIOS", RutaAudiosPorDefecto)
}

/**
 * @brief Obtiene la URL del servidor de metadatos (variable URL_METADATOS).
 * @return URL base sin barra final.
 */
func ObtenerURLMetadatos() string {
	return strings.TrimRight(obtenerVariable("URL_METADATOS", URLMetadatosPorDefecto), "/")
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
