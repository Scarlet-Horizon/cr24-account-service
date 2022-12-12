package messaging

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Messaging struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func (receiver Messaging) Init() error {
	conn, err := amqp.Dial("ampq://studentdocker.informatika.uni-mb.si:5672")
	if err != nil {
		return err
	}
	receiver.conn = conn

	ch, err := receiver.conn.Channel()
	if err != nil {
		return err
	}
	receiver.channel = ch

	q, err := ch.QueueDeclare(
		"account-service",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	receiver.queue = q
	return nil
}

func (receiver Messaging) Close() {
	if receiver.conn != nil {
		err := receiver.conn.Close()
		log.Printf("conn close error: %v", err)
	}
	if receiver.channel != nil {
		err := receiver.channel.Close()
		log.Printf("channel close error: %v", err)
	}
}
