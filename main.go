package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"go-stack-app/questions"
	"go-stack-app/settings"
)

var wg = sync.WaitGroup{}

func sendQuestionToQueue(question questions.Item) {
	//This is mock
	fmt.Printf("Sending question with id: %v to queue...\n", question.QuestionID)
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
				go sendQuestionToQueue(item)
			}
			hasMore = result.HasMore
			page++
			time.Sleep(time.Duration(delay) * time.Millisecond)
			log.Printf("Add delay, to fit in request rate limit. Delayed %v milliseconds...", delay)
			log.Println("Has more pages?", hasMore)
			if newFromDate < result.GetLatesDate() {
				newFromDate = result.GetLatesDate()
			}
		}
		wg.Wait()
		log.Println("Done. sleep for 5 minutes")
		fromDate = newFromDate
		log.Println("New 'fromDate' is now: ", fromDate)
		time.Sleep(5 * time.Minute)
	}
}
