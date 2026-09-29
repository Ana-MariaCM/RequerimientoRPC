/**
 * @file rabbitmqConsumidor.go
 * @brief Componente que consume los mensajes de la cola de reproducciones en RabbitMQ.
 */
package componenteconexioncola

import (
	"fmt"

	"github.com/streadway/amqp"
)

/**
 * @brief Función que procesa el cuerpo de un mensaje consumido.
 */
type ManejadorMensaje func(cuerpo []byte) error

/**
 * @brief Consumidor de mensajes de una cola de RabbitMQ.
 */
type RabbitConsumidor struct {
	conexion   *amqp.Connection ///< Conexión con el broker.
	canal      *amqp.Channel    ///< Canal AMQP abierto sobre la conexión.
	nombreCola string           ///< Nombre de la cola de la que se consume.
}

/**
 * @brief Conecta con RabbitMQ y declara la cola (durable) de la que se consumirá.
 * @param url URL AMQP del broker.
 * @param nombreCola Nombre de la cola.
 * @return El consumidor conectado o un error si no fue posible conectarse.
 */
func NuevoRabbitConsumidor(url string, nombreCola string) (*RabbitConsumidor, error) {
	conexion, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("error conectando a RabbitMQ: %v", err)
	}

	canal, err := conexion.Channel()
	if err != nil {
		conexion.Close()
		return nil, fmt.Errorf("error abriendo canal: %v", err)
	}

	_, err = canal.QueueDeclare(nombreCola, true, false, false, false, nil)
	if err != nil {
		canal.Close()
		conexion.Close()
		return nil, fmt.Errorf("error declarando cola: %v", err)
	}

	// Se entrega un mensaje a la vez para no sobrecargar el servidor.
	if err := canal.Qos(1, 0, false); err != nil {
		canal.Close()
		conexion.Close()
		return nil, fmt.Errorf("error configurando QoS: %v", err)
	}

	return &RabbitConsumidor{conexion: conexion, canal: canal, nombreCola: nombreCola}, nil
}

/**
 * @brief Consume mensajes de la cola indefinidamente.
 *
 * Por cada mensaje imprime un eco, invoca el manejador y confirma (ack) el
 * mensaje. La función retorna cuando se cierra la conexión con el broker.
 * @param manejador Función que procesa el cuerpo de cada mensaje.
 * @return Error si no se pudo registrar el consumidor.
 */
func (thisC *RabbitConsumidor) Consumir(manejador ManejadorMensaje) error {
	mensajes, err := thisC.canal.Consume(
		thisC.nombreCola, // cola
		"",               // consumidor
		false,            // autoAck (se confirma manualmente)
		false,            // exclusive
		false,            // noLocal
		false,            // noWait
		nil,              // argumentos
	)
	if err != nil {
		return fmt.Errorf("error registrando el consumidor: %v", err)
	}

	for mensaje := range mensajes {
		fmt.Printf("\n[Cola] Mensaje consumido de la cola \"%s\": %s\n", thisC.nombreCola, string(mensaje.Body))

		if err := manejador(mensaje.Body); err != nil {
			fmt.Println("[Cola] Mensaje descartado:", err)
			mensaje.Nack(false, false)
			continue
		}
		mensaje.Ack(false)
	}
	return nil
}

/**
 * @brief Cierra el canal y la conexión con RabbitMQ.
 */
func (thisC *RabbitConsumidor) Cerrar() {
	thisC.canal.Close()
	thisC.conexion.Close()
}
