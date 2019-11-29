package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	tb "gopkg.in/tucnak/telebot.v2"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancel := context.WithCancel(context.Background())
	dbCancel = cancel
	initDB(ctx)
}

func main() {
	defer dbCancel()

	pref := getBotPreferences()
	b, err := tb.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	setHandlers(b)
	log.Println("Bot Started!")
	b.Start()
}
