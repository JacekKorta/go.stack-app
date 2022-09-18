package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"sync"

	"go-stack-app/questions"
	"go-stack-app/settings"
)
var wg = sync.WaitGroup{}
	
func sendQuestionToQueue(question questions.Item) {
	//This is mock
	fmt.Printf("Sending %v to queue...\n", question.Title)
	wg.Done()
}



func main() {

	settings := &settings.Settings{}

	settings.GetSettings()

	myClient := &http.Client{Timeout: 10 * time.Second}
	
	page := 26
	fromDate := 0
	hasMore := true
	errorsCount := 0

	for hasMore {
		if errorsCount > 2 {
			log.Println("There is some constant problem with API. Closing program...")
			break
		}
		questionsClient := questions.NewClient(myClient)
		result, err := questionsClient.GetQuestions(settings, page, fromDate)

		if err != nil {
			log.Println("There is so problem with stackoverflow API")
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

	}

	wg.Wait()
}