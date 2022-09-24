package settings

import (
	"fmt"
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
	Rabbit RabbitmQSettings
}

type RabbitmQSettings struct {
	User string
	Password string
	Host string
	Port int
}

func (s *Settings) GetSettings() *Settings {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file")
	}
	s.AppUrl = os.Getenv("APP_URL")
	s.Filter = os.Getenv("FILTER")
	s.Tagged = os.Getenv("TAGGED")
	s.RequestLimit = s.StrToIntParseOrGetDefault("REQEST_LIMIT_PER_SEC", 10)
	s.CheckDelay = s.StrToIntParseOrGetDefault("DELAY_BETWEEN_CHECKS", 5) //minutes
	s.Rabbit.User = os.Getenv("RABBITMQ_USER")
	s.Rabbit.Password = os.Getenv("RABBITMQ_PASSWORD")
	s.Rabbit.Host = os.Getenv("RABBITMQ_HOST")
	s.Rabbit.Port = s.StrToIntParseOrGetDefault("RABBITMQ_PORT", 5672)
	
	return s

}
func (s *Settings) GetRabbitmqUrl(vhost string) string {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
	s.Rabbit.User,
	s.Rabbit.Password,
	s.Rabbit.Host,
	s.Rabbit.Port,
	vhost,
	)
	return url
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