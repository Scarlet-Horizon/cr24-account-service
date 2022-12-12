package messaging

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"strings"
	"time"
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

func (receiver Messaging) Info(context *gin.Context) {
	log.Printf("time=%v", time.Now().Format("2006-01-02 15-04-05"))
	log.Printf("id=%s", uuid.NewString())
	log.Printf("level=info")
	log.Printf("path=%s", context.Request.RequestURI)

	correlation := context.GetHeader("Correlation")
	if correlation == "" {
		correlation = "nil"
	}

	log.Printf("correlation=%s", correlation)
	log.Printf("ip=%s\n", context.ClientIP())

	var token string
	values := strings.Split(context.GetHeader("Authorization"), "Bearer ")
	if len(values) == 2 {
		token = values[1]
	} else {
		token = "nil"
	}

	log.Printf("auth=%s", token)
}
