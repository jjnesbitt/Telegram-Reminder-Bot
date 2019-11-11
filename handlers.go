package main

import (
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func remindMeHandler(b *tb.Bot) func(m *tb.Message) {
	return func(m *tb.Message) {
		wait, err := getWaitTime(m.Payload)
		if m.Payload == "" || err {
			b.Send(m.Sender, "No valid match!")
		}

		duration := time.Duration(int64(wait.seconds * int(time.Second)))

		confirmReminderSet(&wait, b, m.Sender)

		time.Sleep(duration)
		b.Send(m.Sender, "You entered "+m.Payload)
	}
}
