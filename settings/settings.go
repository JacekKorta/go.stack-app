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
	CheckDelay int
}

func (s *Settings) GetSettings() *Settings {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file")
	}
	s.AppUrl = os.Getenv("APP_URL")
	s.Filter = os.Getenv("FILTER")
	s.Tagged = os.Getenv("TAGGED")

	strRequestLimit := os.Getenv("REQEST_LIMIT_PER_SEC")
	intRequestLimit, err := strconv.Atoi(strRequestLimit)
	if err != nil {
		s.RequestLimit = 50
	} else {
		s.RequestLimit = intRequestLimit
	}

	strDelay := os.Getenv("DELAY_BETWEEN_CHECKS")
	intDelay, err := strconv.Atoi(strDelay)
	if err != nil {
		s.CheckDelay = 5
	} else {
		s.CheckDelay = intDelay
	}
	
	return s

}

func (s *Settings) GetMilisecondRateLimit() int {
	rest := 1000 % s.RequestLimit
	base := 1000 - rest
	return base / s.RequestLimit
}