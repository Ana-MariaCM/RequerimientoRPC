/**
 * @file rabbitmqPublicador.go
 * @brief Componente que publica mensajes en la cola de RabbitMQ.
 *
 * Permite enviar asincrónicamente la información de las reproducciones al
 * servidor de estadísticas, desacoplando ambos servidores.
 */
package componenteconexioncola

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/streadway/amqp"
)

/**
 * @brief Publicador de mensajes JSON en una cola de RabbitMQ.
 */
type RabbitPublicador struct {
	mu         sync.Mutex       ///< Serializa las publicaciones (el canal AMQP no es concurrente).
	url        string           ///< URL de conexión al broker.
	nombreCola string           ///< Nombre de la cola destino.
	conexion   *amqp.Connection ///< Conexión con el broker.
	canal      *amqp.Channel    ///< Canal AMQP abierto sobre la conexión.
}

/**
 * @brief Crea un publicador y establece la conexión con RabbitMQ.
 * @param url URL AMQP del broker (por ejemplo amqp://guest:guest\@localhost:5672/).
 * @param nombreCola Nombre de la cola en la que se publicarán los mensajes.
 * @return El publicador creado. Si la conexión falla se devuelve también el error;
 *         el publicador intentará reconectarse en la siguiente publicación.
 */
func NuevoRabbitPublicador(url string, nombreCola string) (*RabbitPublicador, error) {
	publicador := &RabbitPublicador{url: url, nombreCola: nombreCola}
	publicador.mu.Lock()
	defer publicador.mu.Unlock()
	return publicador, publicador.conectar()
}

/**
 * @brief Abre la conexión, el canal y declara la cola (durable).
 * @return Error si no fue posible conectarse con el broker.
 */
func (thisP *RabbitPublicador) conectar() error {
	conexion, err := amqp.Dial(thisP.url)
	if err != nil {
		return fmt.Errorf("error conectando a RabbitMQ: %v", err)
	}

	canal, err := conexion.Channel()
	if err != nil {
		conexion.Close()
		return fmt.Errorf("error abriendo canal: %v", err)
	}

	_, err = canal.QueueDeclare(
		thisP.nombreCola, // nombre
		true,             // durable
		false,            // autoDelete
		false,            // exclusive
		false,            // noWait
		nil,              // argumentos
	)
	if err != nil {
		canal.Close()
		conexion.Close()
		return fmt.Errorf("error declarando cola: %v", err)
	}

	thisP.conexion = conexion
	thisP.canal = canal
	return nil
}

/**
 * @brief Serializa un mensaje a JSON y lo publica en la cola.
 *
 * Si no existe una conexión activa se intenta reconectar antes de publicar.
 * @param mensaje Objeto (DTO) a publicar.
 * @return Error si no se pudo serializar o publicar el mensaje.
 */
func (thisP *RabbitPublicador) Publicar(mensaje any) error {
	cuerpo, err := json.Marshal(mensaje)
	if err != nil {
		return fmt.Errorf("error convirtiendo mensaje a JSON: %v", err)
	}

	thisP.mu.Lock()
	defer thisP.mu.Unlock()

	if thisP.conexion == nil || thisP.conexion.IsClosed() {
		if err := thisP.conectar(); err != nil {
			return err
		}
	}

	err = thisP.canal.Publish(
		"",               // exchange por defecto
		thisP.nombreCola, // routing key = nombre de la cola
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         cuerpo,
		},
	)
	if err != nil {
		thisP.conexion = nil
		return fmt.Errorf("error publicando mensaje: %v", err)
	}

	fmt.Printf("[Cola] Mensaje publicado en la cola \"%s\": %s\n", thisP.nombreCola, string(cuerpo))
	return nil
}

/**
 * @brief Cierra el canal y la conexión con RabbitMQ.
 */
func (thisP *RabbitPublicador) Cerrar() {
	thisP.mu.Lock()
	defer thisP.mu.Unlock()

	if thisP.canal != nil {
		thisP.canal.Close()
	}
	if thisP.conexion != nil {
		thisP.conexion.Close()
	}
}
