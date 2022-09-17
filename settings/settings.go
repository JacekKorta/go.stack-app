package settings

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Settings struct {
	AppUrl string
	Filter string
	Tagged string
}

func (s *Settings) GetSettings() *Settings {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file")
	}
	s.AppUrl = os.Getenv("APP_URL")
	s.Filter = os.Getenv("FILTER")
	s.Tagged = os.Getenv("TAGGED")
	return s

}