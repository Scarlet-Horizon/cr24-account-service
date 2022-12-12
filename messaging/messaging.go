package messaging

import (
	"context"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"main/util"
	"os"
	"time"
)

type Messaging struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func (receiver Messaging) Init() error {
	conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
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

func (receiver Messaging) write(message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return receiver.channel.PublishWithContext(ctx,
		"SIPIA-4",
		receiver.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

func (receiver Messaging) WriteInfo(context *gin.Context) {
	err := receiver.write(util.Info(context))
	if err != nil {
		log.Printf("error with messaging info: %s\n", err)
	}
}

func (receiver Messaging) WriteError(errDesc string, context *gin.Context)  {
	err := receiver.write(util.Error(errDesc, context))
	if err != nil {
		log.Printf("error with messaging error: %s\n", err)
	}
}