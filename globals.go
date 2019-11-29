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
	envMode               string
	publicURL             string
	rootTelegramMethodURL string
	mongoURL              string
	mongoPort             string
	currentLimboUsers     map[int]*tb.Message = make(map[int]*tb.Message)
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	envBotToken = os.Getenv("BOT_TOKEN")
	rootTelegramMethodURL = "https://api.telegram.org/bot" + envBotToken

	envRootPublicURL = os.Getenv("ROOT_PUBLIC_URL")
	publicURL = envRootPublicURL + "/" + envBotToken

	envListenPort = os.Getenv("LISTEN_PORT")
	envMode = os.Getenv("BOT_MODE")

	mongoURL = os.Getenv("MONGO_URL")
	mongoPort = os.Getenv("MONGO_PORT")
}

// var dotenvErr = godotenv.Load()

// var envBotToken = os.Getenv("BOT_TOKEN")
// var envListenPort = os.Getenv("LISTEN_PORT")
// var envRootPublicURL = os.Getenv("ROOT_PUBLIC_URL")
// var envMode = os.Getenv("BOT_MODE")

// var publicURL = envRootPublicURL + "/" + envBotToken

// var rootTelegramMethodURL = "https://api.telegram.org/bot" + envBotToken

// var currentLimboUsers map[int]*tb.Message = make(map[int]*tb.Message)
