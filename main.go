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
	
	page := 1
	fromDate := 0

	questionsClient := questions.NewClient(myClient)
	result, err := questionsClient.GetQuestions(settings, page, fromDate)

	if err != nil {
		log.Println("There is so problem with stackoverflow API")
	}
	for _, item := range result.Items {
		wg.Add(1)
		go sendQuestionToQueue(item)
	}
	wg.Wait()
}