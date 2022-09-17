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

//https://api.stackexchange.com/2.3/questions?page=2&pagesize=1&order=desc&max=1663286400&sort=activity&tagged=python&site=stackoverflow&filter=!0WJ3YL2.EQ8B_wPSrO73X35Fv
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

type Client struct {
	http *http.Client
	page int
	tagged string
}

func (c *Client) GetQuestions(settings *settings.Settings, page int, fromDate int) (*QuestionsSearchOut, error) {
	t := time.Now().AddDate(0, -1, 0)
	max := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix() //today in unix
	endpoint := fmt.Sprintf(
		"%s/2.3/questions?page=%d&fromdate=%dorder=desc&max=%d&sort=activity&tagged=%s&site=stackoverflow&filter=%s", 
		settings.AppUrl,
		page,
		fromDate,
		max,
		settings.Tagged,
		settings.Filter,
	)

	log.Println("Making request to: ", endpoint)
	resp, err := c.http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}

	res := &QuestionsSearchOut{}
	return res, json.Unmarshal(body, res)
}