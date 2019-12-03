package main

import (
	"fmt"
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
)

func setHandlers() {
	botInstance.Handle("/remindme", remindMeHandler)
	botInstance.Handle("/cancel", cancelHandler)
	botInstance.Handle(tb.OnText, onTextHandler)
}

func remindMeHandler(m *tb.Message) {
	wait, err := getWaitTime(m.Payload)
	if m.Payload == "" || err != nil {
		botInstance.Send(m.Chat, "No valid time units found!")
		return
	}

	if m.ReplyTo == nil {
		botInstance.Send(m.Chat, "No message to forward!")
		return
	}

	go confirmReminderSet(&wait, m.Chat)
	go forwardMessageAfterDelay(wait, m.Sender, m.ReplyTo)
}

func cancelHandler(m *tb.Message) {
	if !m.Private() {
		botInstance.Send(m.Chat, "Must be done in private chat!")
		return
	}

	// cur, err := dbCol.Find(dbCtx, bson.M{"user.id": m.Sender.ID})
	// defer cur.Close(dbCtx)

	// if err != nil {
	// 	botInstance.Send(m.Chat, "Error Finding Reminders!")
	// }

	// var reminders []StoredReminder
	// cur.All(dbCtx, &reminders)

	reminders, err := getUserReminders(m.Sender)
	if err != nil {
		botInstance.Send(m.Chat, "Error Finding Reminders!")
		return
	}

	if len(reminders) == 0 {
		botInstance.Send(m.Chat, "No Pending Reminders!")
		return
	}

	for i := range reminders {
		fmt.Println(reminders[i].User.Username)

		messageID := reminders[i].MessageID
		chat := tb.Chat{ID: reminders[i].ChatID}
		message := tb.Message{ID: messageID, Chat: &chat}

		forwardMessage(m.Sender, &message)
	}
}

// Handles direct forwarded requests
func onTextHandler(m *tb.Message) {
	if !m.Private() {
		return
	}

	log.Println("Private Message!")

	if waitingMessage, ok := currentLimboUsers[m.Sender.ID]; ok {
		// Already waiting
		wait, err := getWaitTime(m.Text)
		if err != nil {
			botInstance.Send(m.Sender, "No valid match! Aborting...")
		} else {
			go confirmReminderSet(&wait, m.Sender)
			go forwardMessageAfterDelay(wait, m.Sender, waitingMessage)
		}

		delete(currentLimboUsers, m.Sender.ID)
	} else {
		currentLimboUsers[m.Sender.ID] = m
		botInstance.Send(m.Sender, "When should I remind you?")
	}
}
