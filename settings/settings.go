package settings

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Settings struct {
	AppUrl string
	Filter string
	Tagged string
	RequestLimit int
}

func (s *Settings) GetSettings() *Settings {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file")
	}
	s.AppUrl = os.Getenv("APP_URL")
	s.Filter = os.Getenv("FILTER")
	s.Tagged = os.Getenv("TAGGED")
	requestLimitStr := os.Getenv("REQEST_LIMIT_PER_SEC")
	intVar, err := strconv.Atoi(requestLimitStr)
	if err != nil {
		s.RequestLimit = 0
	} else {
		s.RequestLimit = intVar
	}
	
	return s

}

func (s *Settings) GetMilisecondRateLimit() int {
	rest := 1000 % s.RequestLimit
	base := 1000 - rest
	return base / s.RequestLimit
}