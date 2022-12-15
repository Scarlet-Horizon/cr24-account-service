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
	queue   *amqp.Queue
}

func (receiver *Messaging) Init() error {
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
		os.Getenv("EXCHANGE_QUEUE_NAME"),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	receiver.queue = &q
	return nil
}

func (receiver *Messaging) Close() {
	if receiver.conn != nil && !receiver.conn.IsClosed() {
		err := receiver.conn.Close()
		if err != nil {
			log.Printf("conn close error: %v", err)
		}
	}
	if receiver.channel != nil && !receiver.channel.IsClosed() {
		err := receiver.channel.Close()
		if err != nil {
			log.Printf("channel close error: %v", err)
		}
	}
}

func (receiver *Messaging) write(message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return receiver.channel.PublishWithContext(ctx,
		os.Getenv("EXCHANGE_QUEUE_NAME"),
		receiver.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

func (receiver *Messaging) WriteInfo(context *gin.Context) {
	err := receiver.write(util.Info(context))
	if err != nil {
		log.Printf("error with messaging info: %s\n", err)
	}
}

func (receiver *Messaging) WriteError(context *gin.Context) {
	context.Next()

	for _, err := range context.Errors {
		returnerErr := receiver.write(util.Error(err.Error(), context))
		if returnerErr != nil {
			log.Printf("error with messaging error: %s\n", returnerErr)
		}
	}
}
