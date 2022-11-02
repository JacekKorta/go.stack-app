package messages

import (
	"context"
	"log"

	"go-stack-app/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)
	
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
