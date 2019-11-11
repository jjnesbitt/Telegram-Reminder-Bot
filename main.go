package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	tb "gopkg.in/tucnak/telebot.v2"
)

func setHandlers(b *tb.Bot) {
	b.Handle("/remindme", remindMeHandler)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("BOT_TOKEN")
	listenPort := ":" + os.Getenv("LISTEN_PORT")
	publicURL := os.Getenv("PUBLIC_URL")

	pref := tb.Settings{
		Token: botToken,
		Poller: &tb.Webhook{
			Listen:   listenPort,
			Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
		},
	}

	b, err := tb.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	setHandlers(b)
	println("Bot Started!")
	b.Start()
}
