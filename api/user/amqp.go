package user

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"time"
)

type MQ interface {
	PublishMessage(queue string, body string)
}

type mq struct {
	conn *amqp.Connection
}

func NewMQ(conn *amqp.Connection) *mq {
	return &mq{conn}
}

// PublishMessage sends message to RabbitMQ, where body contains user id
// and queueName in [user_create, user_update, user_delete]
// for other services notification about user changes
func (m *mq) PublishMessage(queueName string, body string) {
	ch, err := m.conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "plain/text",
			Body:        []byte(body),
		})

	if err != nil {
		log.Errorln("failed to send message", err)
	} else {
		log.Infoln("message sent", queueName, body)
	}
}
