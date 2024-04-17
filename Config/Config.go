package Config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	MongodbURI           string
	ServerAddress        string
	AccessTokenLifeSpan  string
	RefreshTokenLifeSpan string
}

func Init() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		MongodbURI:           os.Getenv("MONGOURI"),
		ServerAddress:        os.Getenv("SERVERADDRESS"),
		AccessTokenLifeSpan:  os.Getenv("ACCESS_TOKEN_LIFESPAN"),
		RefreshTokenLifeSpan: os.Getenv("REFRESH_TOKEN_LIFESPAN"),
	}
}
