package main

import (
	"log"

	"github.com/joho/godotenv"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pref := getBotPreferences()
	b, err := tb.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	setHandlers(b)
	log.Println("Bot Started!")
	b.Start()
}
