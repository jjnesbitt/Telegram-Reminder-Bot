package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	tb "gopkg.in/tucnak/telebot.v2"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	body, _ := r.GetBody()
	fmt.Printf("<%s>", body)
}

// func setHandlers(b *tb.Bot, s *scheduler.Scheduler) {
// 	b.Handle("/remindme", func(m *tb.Message) {
// 		b.Send(m.Sender, "You entered"+m.Payload)
// 	})
// }

func setHandlers(b *tb.Bot) {
	b.Handle("/remindme", func(m *tb.Message) {
		b.Send(m.Sender, "You entered"+m.Payload)
	})
}

func main() {
	// http.HandleFunc("/", handler)
	// fmt.Printf("Listening!\n")
	// http.ListenAndServe(":8080", nil)

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

	// sqliteStorage := storage.NewSqlite3Storage()
	// if err := storage.Connect(); err != nil {
	// 	log.Fatal("Could not connect to db", err)
	// }
	// if err := storage.Initialize(); err != nil {
	// 	log.Fatal("Could not intialize database", err)
	// }

	// s := scheduler.New(sqliteStorage)
	setHandlers(b)
	go b.Start()
	println("Bot Started!")
}
