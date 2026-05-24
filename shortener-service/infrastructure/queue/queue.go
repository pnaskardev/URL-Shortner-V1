package queue

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewQueueClient() (*QueueClient, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &QueueClient{
		conn:    conn,
		channel: ch,
	}, nil
}

func (c *QueueClient) DeclareQueue(name string) error {
	_, err := c.channel.QueueDeclare(
		// "hello", // name
		name,
		true,  // durability
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		amqp.Table{
			amqp.QueueTypeArg: amqp.QueueTypeQuorum,
		},
	)
	return err
}

func (c *QueueClient) Publish(queueName string, body []byte) error {
	return c.channel.Publish(
		"url_creation",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}
