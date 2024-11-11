package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var ENV struct {
	PORT   string
	DB_URL string
	SECRET string
}

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ENV.PORT = os.Getenv("PORT")
	ENV.DB_URL = os.Getenv("DB_URL")
	ENV.SECRET = os.Getenv("SECRET")

}

var UserString = "user"
