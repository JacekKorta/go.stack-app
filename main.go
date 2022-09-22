package main

import (
	"log"
	"net/http"
	"sync"
	"time"
	"context"

	"go-stack-app/questions"
	"go-stack-app/settings"
	amqp "github.com/rabbitmq/amqp091-go"
)



func failOnError(err error, msg string) {
if err != nil {
  log.Panicf("%s: %s", msg, err)
}
}
var wg = sync.WaitGroup{}

func publishMessage(ctx context.Context, body string, ch *amqp.Channel) {
	err := ch.PublishWithContext(ctx,
	"stack.questions.raw",     // exchange
	"stack.questions.duplicated", // routing key
	false,  // mandatory
	false,  // immediate
	amqp.Publishing {
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
	wg.Done()
}


func main() {

	settings := &settings.Settings{}

	settings.GetSettings()

	myClient := &http.Client{Timeout: 10 * time.Second}

	page := 1
	fromDate := 0
	hasMore := true
	errorsCount := 0
	maxErrorCount := 2
	delay := settings.GetMilisecondRateLimit()
	var newFromDate int = 0

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/mtg")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		for hasMore {
			if errorsCount > 2 {
				log.Println("There is some constant problem with API. Breaking current loop...")
				//Todo: change to sleep for a while, and remove break
				break
			}
			questionsClient := questions.NewClient(myClient)
			result, err := questionsClient.GetQuestions(settings, page, fromDate)

			if err != nil {
				log.Printf(
					"There is so problem with stackoverflow API. Error count: %v/%v\n", errorsCount, maxErrorCount,
				)

				errorsCount++
				continue
			}
			errorsCount = 0
			for _, item := range result.Items {
				wg.Add(1)
				go publishMessage(ctx, item.Title, ch)
			}
			hasMore = result.HasMore
			page++
			time.Sleep(time.Duration(delay) * time.Millisecond)
			log.Printf("Add delay, to fit in request rate limit. Delayed %v milliseconds...", delay)
			log.Println("Has more pages?", hasMore)
			if newFromDate < result.GetLatesDate() {
				newFromDate = result.GetLatesDate()
			}
			wg.Wait()
		}
		
		log.Println("Done. sleep for 5 minutes")
		fromDate = newFromDate
		log.Println("New 'fromDate' is now: ", fromDate)
		time.Sleep(5 * time.Minute)
		errorsCount = 0
		page = 1
		hasMore = true

	}
}
