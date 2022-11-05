package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go-stack-app/questions"
	"go-stack-app/settings"
	"go-stack-app/messages"
	"go-stack-app/utils"

)


func main() {

	settings := &settings.Settings{}

	settings.GetSettings()

	myClient := &http.Client{Timeout: 10 * time.Second}

	page := 1
	fromDate := 0
	hasMore := true
	errorsCount := 0
	maxErrorCount := 2
	sleepAfterGrab := settings.CheckDelay
	delay := settings.GetMilisecondRateLimit()
	var newFromDate int = 0

	conn := messages.GetConnection(settings.GetRabbitmqUrl("/mtg"))
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		for hasMore {
			if errorsCount > 2 {
				log.Println("There is some constant problem with API. Breaking current loop...")
				//Todo: change to sleep for a while, and remove break?
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
				body, err := json.Marshal(item)
				utils.FailOnError(err, "Unable to marshal item")	
				go messages.PublishMessage(ctx, string(body), ch, item.QuestionID)
			}
			hasMore = result.HasMore
			page++
			time.Sleep(time.Duration(delay) * time.Millisecond)
			log.Printf("Add delay, to fit in request rate limit. Delayed %v milliseconds...", delay)
			log.Printf("Has more pages? %v\n", hasMore)
			if newFromDate < result.GetLatesDate() {
				newFromDate = result.GetLatesDate()
			}
		}
		
		log.Printf("Done. sleep for %d minutes\n", sleepAfterGrab)
		fromDate = newFromDate
		log.Printf("New 'fromDate' is now: %v\n", fromDate)
		time.Sleep(time.Duration(sleepAfterGrab) * time.Minute)
		errorsCount = 0
		page = 1
		hasMore = true

	}
}
