package main

import (
	// "encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	// "io"
	// "net/http"

	"go-stack-app/questions"
	"go-stack-app/settings"
)
	


func main() {

	settings := &settings.Settings{}

	settings.GetSettings()

	myClient := &http.Client{Timeout: 10 * time.Second}
	
	page := 1
	fromDate := 0

	questionsClient := questions.NewClient(myClient)
	result, err := questionsClient.GetQuestions(settings, page, fromDate)

	if err != nil {
		log.Println("Api nie działa")
	}
	fmt.Println(result)
	fmt.Println(result.GetLatesDate())
}