package componenteconexioncola

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

type NotificacionCancion struct {
	Titulo  string `json:"titulo"`
	Artista string `json:"artista"`
	Genero  string `json:"genero"`
	Mensaje string `json:"mensaje"`
}

func NewRabbitPublisher() (*RabbitPublisher, error) {
	conn, err := amqp.Dial("amqp://admin:1234@10.113.184.136:5672")

	if err != nil {
		return nil, fmt.Errorf("Error conectando a RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Error abriendo canal: %v", err)
	}

	q, err := ch.QueueDeclare(
		"notificaciones_canciones",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("Error declarando cola: %v", err)
	}

	return &RabbitPublisher{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

func (p *RabbitPublisher) PublicarNotificacion(msg NotificacionCancion) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("Error convirtiendo mensaje a json: %v", err)
	}
	err = p.channel.Publish(
		"",
		p.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "aplication/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("Error publicando mensaje: %v", err)
	}
	fmt.Println("Notificacion enviada a RabbitMQ:  ", string(body))
	return nil
}

func (p *RabbitPublisher) Cerrar() {
	p.channel.Close()
	p.conn.Close()
}
