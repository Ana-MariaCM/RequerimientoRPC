/**
 * @file configuracion.go
 * @brief Parámetros de configuración del servidor de estadísticas.
 */
package configuracion

import "os"

/** @brief URL por defecto del broker RabbitMQ. */
const URLRabbitMQPorDefecto = "amqp://guest:guest@localhost:5672/"

/** @brief Nombre de la cola de la que se consumen las reproducciones. */
const NombreColaReproducciones = "reproducciones_audios"

/**
 * @brief Obtiene la URL de conexión a RabbitMQ (variable RABBITMQ_URL).
 * @return URL AMQP del broker.
 */
func ObtenerURLRabbitMQ() string {
	valor := os.Getenv("RABBITMQ_URL")
	if valor == "" {
		return URLRabbitMQPorDefecto
	}
	return valor
}
