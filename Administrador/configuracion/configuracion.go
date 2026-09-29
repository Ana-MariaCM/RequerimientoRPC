/**
 * @file configuracion.go
 * @brief Parámetros de configuración del administrador.
 */
package configuracion

import (
	"os"
	"strings"
)

/** @brief URL por defecto del servidor de audios (REST). */
const URLServidorAudiosPorDefecto = "http://localhost:5000"

/**
 * @brief Obtiene la URL base del servidor de audios (variable URL_SERVIDOR_AUDIOS).
 * @return URL base sin barra final.
 */
func ObtenerURLServidorAudios() string {
	valor := os.Getenv("URL_SERVIDOR_AUDIOS")
	if valor == "" {
		valor = URLServidorAudiosPorDefecto
	}
	return strings.TrimRight(valor, "/")
}
