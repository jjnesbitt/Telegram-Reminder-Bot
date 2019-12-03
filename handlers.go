package main

import (
	"fmt"
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
)

func setHandlers(b *tb.Bot) {
	b.Handle("/remindme", remindMeHandler(b))
	b.Handle("/cancel", cancelHandler(b))
	b.Handle(tb.OnText, onTextHandler(b))
}

func remindMeHandler(b *tb.Bot) func(m *tb.Message) {
	return func(m *tb.Message) {
		wait, err := getWaitTime(m.Payload)
		if m.Payload == "" || err != nil {
			b.Send(m.Chat, "No valid time units found!")
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

func cancelHandler(b *tb.Bot) func(m *tb.Message) {
	return func(m *tb.Message) {
		if !m.Private() {
			b.Send(m.Chat, "Must be done in private chat!")
			return
		}

		// cur, err := dbCol.Find(dbCtx, bson.M{"user.id": m.Sender.ID})
		// defer cur.Close(dbCtx)

		// if err != nil {
		// 	b.Send(m.Chat, "Error Finding Reminders!")
		// }

		// var reminders []StoredReminder
		// cur.All(dbCtx, &reminders)

		reminders, err := getUserReminders(m.Sender)
		if err != nil {
			b.Send(m.Chat, "Error Finding Reminders!")
			return
		}

		if len(reminders) == 0 {
			b.Send(m.Chat, "No Pending Reminders!")
			return
		}

		for i := range reminders {
			fmt.Println(reminders[i].User.Username)

			messageID := reminders[i].MessageID
			chat := tb.Chat{ID: reminders[i].ChatID}
			message := tb.Message{ID: messageID, Chat: &chat}

			forwardMessage(b, m.Sender, &message)
		}
	}
}

// Handles direct forwarded requests
func onTextHandler(b *tb.Bot) func(m *tb.Message) {
	return func(m *tb.Message) {
		if !m.Private() {
			return
		}

		log.Println("Private Message!")

		if waitingMessage, ok := currentLimboUsers[m.Sender.ID]; ok {
			// Already waiting
			wait, err := getWaitTime(m.Text)
			if err != nil {
				b.Send(m.Sender, "No valid match! Aborting...")
			} else {
				go confirmReminderSet(&wait, b, m.Sender)
				go forwardMessageAfterDelay(wait.duration, b, m.Sender, waitingMessage)
			}

			delete(currentLimboUsers, m.Sender.ID)
		} else {
			currentLimboUsers[m.Sender.ID] = m
			b.Send(m.Sender, "When should I remind you?")
		}
	}
}
