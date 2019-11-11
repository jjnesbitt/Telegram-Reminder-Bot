package main

import (
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func remindMeHandler(b *tb.Bot) func(m *tb.Message) {
	return func(m *tb.Message) {
		if m.Sender.IsBot {
			return
		}

		if m.IsForwarded() {
			// Reply with keyboard for setting timer
			return
		}

		wait, err := getWaitTime(m.Payload)
		if m.Payload == "" || err {
			b.Send(m.Sender, "No valid match!")
			return
		}

		duration := time.Duration(int64(wait.seconds * int(time.Second)))

		if m.ReplyTo == nil {
			b.Send(m.Chat, "No message to forward!")
			return
		}

		confirmReminderSet(&wait, b, m.Sender)
		// go sendMessage(duration, m.ReplyTo, b)
		time.Sleep(duration)
		b.Forward(m.Sender, m.ReplyTo)
	}
}
