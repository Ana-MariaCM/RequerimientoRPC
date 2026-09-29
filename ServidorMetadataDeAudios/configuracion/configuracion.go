/**
 * @file configuracion.go
 * @brief Parámetros de configuración del servidor de metadatos de audios.
 *
 * Los valores por defecto pueden sobrescribirse mediante variables de entorno,
 * lo que permite desplegar el servidor en una máquina distinta sin recompilar.
 */
package configuracion

import "os"

/** @brief Puerto HTTP por defecto en el que escucha el servidor REST de metadatos. */
const PuertoPorDefecto = "5001"

/**
 * @brief Archivo JSON por defecto donde se guardan los metadatos de los audios
 *        que registra el administrador (ruta relativa a la carpeta del servidor).
 */
const RutaAudiosRegistradosPorDefecto = "datos/audiosRegistrados.json"

/**
 * @brief Obtiene la ruta del archivo JSON de audios registrados (variable RUTA_AUDIOS_REGISTRADOS).
 * @return Ruta del archivo JSON.
 */
func ObtenerRutaAudiosRegistrados() string {
	return obtenerVariable("RUTA_AUDIOS_REGISTRADOS", RutaAudiosRegistradosPorDefecto)
}

/**
 * @brief Obtiene el puerto en el que escuchará el servidor REST.
 *
 * Se consulta la variable de entorno PUERTO_METADATOS; si no existe se usa
 * PuertoPorDefecto.
 * @return Puerto en formato texto (por ejemplo "5001").
 */
func ObtenerPuerto() string {
	return obtenerVariable("PUERTO_METADATOS", PuertoPorDefecto)
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
