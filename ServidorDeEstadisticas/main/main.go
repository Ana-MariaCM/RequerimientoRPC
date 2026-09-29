/**
 * @file main.go
 * @brief Punto de entrada del servidor de estadísticas (consumidor de la cola de mensajes).
 */
package main

import (
	"fmt"
	"time"

	controlador "estadisticas/capaControladores"
	componenteconexioncola "estadisticas/componenteConexionCola"
	"estadisticas/configuracion"
)

/** @brief Tiempo de espera entre intentos de conexión con RabbitMQ. */
const esperaReconexion = 5 * time.Second

/**
 * @brief Conecta con RabbitMQ y consume las reproducciones; si la conexión se
 *        pierde, reintenta cada cinco segundos.
 */
func main() {
	ctrl := controlador.NuevoControladorEstadisticas()
	url := configuracion.ObtenerURLRabbitMQ()
	nombreCola := configuracion.NombreColaReproducciones

	for {
		consumidor, err := componenteconexioncola.NuevoRabbitConsumidor(url, nombreCola)
		if err != nil {
			fmt.Println(err, "- reintentando en", esperaReconexion)
			time.Sleep(esperaReconexion)
			continue
		}

		fmt.Printf("Servidor de estadísticas escuchando la cola \"%s\"...\n", nombreCola)
		if err := consumidor.Consumir(ctrl.ProcesarMensaje); err != nil {
			fmt.Println(err)
		}
		consumidor.Cerrar()
		fmt.Println("Conexión con RabbitMQ perdida - reintentando en", esperaReconexion)
		time.Sleep(esperaReconexion)
	}
}
