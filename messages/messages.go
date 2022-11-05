package messages

import (
	"context"
	"log"
	"time"

	"go-stack-app/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

func GetConnection(url string) (*amqp.Connection) {
	var isConneted bool = false
	var connnection amqp.Connection

	for !isConneted {
		conn, err := amqp.Dial(url)
		if err == nil {
			isConneted = true
			connnection = *conn
			break
		}
		log.Println("Can't connect to rabbitmq server. I will try later")
		time.Sleep(30 * time.Second)
		log.Println("Reconect...")
	}
	log.Println("Connecting to RabbitMQ - Done")
	return &connnection
}
	
func PublishMessage(ctx context.Context, body string, ch *amqp.Channel, mark int) {
	err := ch.PublishWithContext(ctx,
	"stack.questions.raw",     // exchange
	"stack.questions.duplicated", // routing key
	false,  // mandatory
	false,  // immediate
	amqp.Publishing {
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	utils.FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent question with id: %v\n", mark)
}
