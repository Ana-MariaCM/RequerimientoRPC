/**
 * @file configuracion.go
 * @brief Parámetros de configuración del servidor de streaming.
 *
 * Todos los valores pueden sobrescribirse mediante variables de entorno.
 */
package configuracion

import "os"

/** @brief Puerto por defecto del servidor gRPC de streaming. */
const PuertoPorDefecto = "50051"

/**
 * @brief Carpeta por defecto donde se encuentran los archivos mp3.
 *
 * Corresponde a la carpeta en la que el servidor de audios almacena los audios
 * que recibe del administrador (recurso Audio.mp3 compartido del diagrama).
 */
const RutaAudiosPorDefecto = "../ServidorDeAudios/audios"

/** @brief URL por defecto del broker RabbitMQ. */
const URLRabbitMQPorDefecto = "amqp://guest:guest@localhost:5672/"

/** @brief Nombre de la cola donde se publican las reproducciones. */
const NombreColaReproducciones = "reproducciones_audios"

/** @brief Tamaño en bytes de cada fragmento enviado por streaming (32 KB). */
const TamanioFragmento = 32 * 1024

/**
 * @brief Obtiene el puerto del servidor gRPC (variable PUERTO_STREAMING).
 * @return Puerto en formato texto.
 */
func ObtenerPuerto() string {
	return obtenerVariable("PUERTO_STREAMING", PuertoPorDefecto)
}

/**
 * @brief Obtiene la carpeta de los archivos mp3 (variable RUTA_AUDIOS).
 * @return Ruta de la carpeta de audios.
 */
func ObtenerRutaAudios() string {
	return obtenerVariable("RUTA_AUDIOS", RutaAudiosPorDefecto)
}

/**
 * @brief Obtiene la URL de conexión a RabbitMQ (variable RABBITMQ_URL).
 * @return URL AMQP del broker.
 */
func ObtenerURLRabbitMQ() string {
	return obtenerVariable("RABBITMQ_URL", URLRabbitMQPorDefecto)
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
