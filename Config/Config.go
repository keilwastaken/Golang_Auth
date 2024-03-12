package Config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	MongodbURI    string
	ServerAddress string
}

func Init() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		MongodbURI:    os.Getenv("MONGOURI"),
		ServerAddress: os.Getenv("SERVERADDRESS"),
	}
}
