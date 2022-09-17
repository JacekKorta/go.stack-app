package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go-stack-app/questions"
	"go-stack-app/settings"
)
	


func main() {

	settings := &settings.Settings{}

	settings.GetSettings()

	base_url := settings.AppUrl

	url := "https://api.stackexchange.com/2.3/questions?page=2&pagesize=1&order=desc&max=1663286400&sort=activity&tagged=python&site=stackoverflow&filter=!)rtsWHVtIIZq9Fi.wGEU"

	res := &questions.QuestionsSearchOut{}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Println(string(body))
	}

	fmt.Println(string(body))
	// res := &ResponseJson{}

	err = json.Unmarshal(body, res)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.Items)
	fmt.Println(base_url)
}