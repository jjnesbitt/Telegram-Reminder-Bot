package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	envBotToken           string
	envListenPort         string
	envRootPublicURL      string
	botMode               string
	publicURL             string
	rootTelegramMethodURL string
	mongoHostname         string
	mongoPort             string
	currentLimboUsers     map[int]*tb.Message = make(map[int]*tb.Message)
)

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func init() {
	err := godotenv.Overload(".env.default", ".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	envBotToken = os.Getenv("BOT_TOKEN")
	rootTelegramMethodURL = "https://api.telegram.org/bot" + envBotToken

	envRootPublicURL = os.Getenv("ROOT_PUBLIC_URL")
	publicURL = envRootPublicURL + "/" + envBotToken

	envListenPort = getEnv("LISTEN_PORT", "9000")
	botMode = getEnv("BOT_MODE", "development")

	mongoHostname = getEnv("MONGO_HOSTNAME", "localhost")
	mongoPort = getEnv("MONGO_PORT", "27017")
}
