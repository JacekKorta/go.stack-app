package questions

import (
	"encoding/json"
	"fmt"
	"go-stack-app/settings"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type QuestionsSearchOut struct {
		Items []Item `json:"items"`
		HasMore        bool `json:"has_more"`
		QuotaMax       int  `json:"quota_max"`
		QuotaRemaining int  `json:"quota_remaining"`
}


type Item struct {
	Tags  []string `json:"tags"`
	IsAnswered       bool   `json:"is_answered"`
	LastActivityDate int    `json:"last_activity_date"`
	CreationDate     int    `json:"creation_date"`
	QuestionID       int    `json:"question_id"`
	Link             string `json:"link"`
	Title            string `json:"title"`
	Body             string `json:"body"`
}


func (q *QuestionsSearchOut) GetLatesDate() int {
	var latestDate int = 0
	for _, item := range(q.Items) {
		if item.CreationDate > latestDate {
			latestDate = item.CreationDate
		}
	}
	return latestDate
}

type QuestionErrorResponse struct {
	ErrorID      int    `json:"error_id"`
	ErrorMessage string `json:"error_message"`
	ErrorName    string `json:"error_name"`
}

type Client struct {
	http *http.Client
}

func (c *Client) GetQuestions(settings *settings.Settings, page int, fromDate int) (*QuestionsSearchOut, error) {
	t := time.Now()
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	if fromDate == 0 {
		fromDate = int(today.AddDate(0, -1, 0).Unix())
	}

	endpoint := fmt.Sprintf(
		"%s/2.3/questions?page=%d&fromdate=%d&order=desc&sort=activity&tagged=%s&site=stackoverflow&filter=%s", 
		settings.AppUrl,
		page,
		fromDate,
		settings.Tagged,
		settings.Filter,
	)
	
	res := &QuestionsSearchOut{}
	log.Printf("Making request to: %v\n", endpoint)
	resp, err := c.http.Get(endpoint)
	if err != nil {
		log.Println(err)
		return res, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return res, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Status code: %v\n", resp.StatusCode)
		log.Printf("Response: %v\n", string(body))
		return res, fmt.Errorf(string(body))
	}

	return res, json.Unmarshal(body, res)
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{httpClient}
}
