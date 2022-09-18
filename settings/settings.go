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
		log.Fatal("Error loading .env file")
	}
	s.AppUrl = os.Getenv("APP_URL")
	s.Filter = os.Getenv("FILTER")
	s.Tagged = os.Getenv("TAGGED")
	s.RequestLimit = s.StrToIntParseOrGetDefault("REQEST_LIMIT_PER_SEC", 50)
	s.CheckDelay = s.StrToIntParseOrGetDefault("DELAY_BETWEEN_CHECKS", 5)
	
	return s

}

func (s *Settings) GetMilisecondRateLimit() int {
	rest := 1000 % s.RequestLimit
	base := 1000 - rest
	return base / s.RequestLimit
}

func (s *Settings) StrToIntParseOrGetDefault(envName string, defaultValue int) int {
	// Method parse env variable to int. If its not posible or env is not set it returns default int value
	strValue := os.Getenv(envName)
	if strValue == "" {
		log.Println("Variable is empty. Using default.")
		return defaultValue
	}
	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		log.Println("Can't use variable. Using default.")
		return defaultValue
	} 
	return intValue
}