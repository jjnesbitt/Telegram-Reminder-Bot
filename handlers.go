package main

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

func setHandlers(b *tb.Bot) {
	b.Handle("/remindme", remindMeHandler(b))
	b.Handle(tb.OnText, onTextHandler(b))
}

func remindMeHandler(b *tb.Bot) func(m *tb.Message) {
	return func(m *tb.Message) {
		if m.Sender.IsBot {
			return
		}

		wait, err := getWaitTime(m.Payload)
		if m.Payload == "" || err {
			b.Send(m.Chat, "No valid match!")
			return
		}

		if m.ReplyTo == nil {
			b.Send(m.Chat, "No message to forward!")
			return
		}

		go confirmReminderSet(&wait, b, m.Chat)
		go forwardMessageAfterDelay(wait.duration, b, m.Sender, m.ReplyTo)
	}
}

// Handles direct forwarded requests
func onTextHandler(b *tb.Bot) func(m *tb.Message) {
	return func(m *tb.Message) {
		fmt.Println(currentLimboUsers)

		if m.Sender.IsBot || (!m.Private() && m.Payload != "") {
			return
		}

		if waitingMessage, ok := currentLimboUsers[m.Sender.ID]; ok {
			// Already waiting
			wait, err := getWaitTime(m.Text)
			if err {
				b.Send(m.Sender, "No valid match! Aborting...")
			} else {
				go confirmReminderSet(&wait, b, m.Sender)
				go forwardMessageAfterDelay(wait.duration, b, m.Sender, waitingMessage)
			}

			delete(currentLimboUsers, m.Sender.ID)
		} else {
			currentLimboUsers[m.Sender.ID] = m
			// options := tb.SendOptions{ReplyMarkup: &numKeyboardMarkup, ReplyTo: m}
			// b.Send(m.Sender, "When should I remind you?", &options)
			b.Send(m.Sender, "When should I remind you?")
		}
	}
}
