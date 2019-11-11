package main

import (
	"os"

	"github.com/joho/godotenv"
	tb "gopkg.in/tucnak/telebot.v2"
)

var dotenvErr = godotenv.Load()

var envBotToken = os.Getenv("BOT_TOKEN")
var envListenPort = os.Getenv("LISTEN_PORT")
var envRootPublicURL = os.Getenv("ROOT_PUBLIC_URL")
var envMode = os.Getenv("BOT_MODE")

var publicURL = envRootPublicURL + "/" + envBotToken

var rootTelegramMethodURL = "https://api.telegram.org/bot" + envBotToken

var currentLimboUsers map[int]*tb.Message = make(map[int]*tb.Message)
