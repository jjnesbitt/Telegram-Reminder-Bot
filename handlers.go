package main

import (
	"fmt"
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func setHandlers() {
	botInstance.Handle("/remindme", remindMeHandler)
	botInstance.Handle("/cancel", cancelHandler)

	botInstance.Handle(tb.OnText, onTextHandler)
	botInstance.Handle(tb.OnCallback, deleteReminderHandler)
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

	go confirmReminderSet(wait, m.Chat)
	go forwardMessageAfterDelay(wait, m.Sender, m.ReplyTo)
}

func cancelHandler(m *tb.Message) {
	if !m.Private() {
		botInstance.Send(m.Chat, "Must be done in private chat!")
		return
	}

	reminders, err := getUserReminders(m.Sender)
	if err != nil {
		botInstance.Send(m.Chat, "Error Finding Reminders!")
		return
	}

	if len(reminders) == 0 {
		botInstance.Send(m.Chat, "No Pending Reminders!")
		return
	}

	var buttonArray [][]tb.InlineButton

	for i := range reminders {
		messageText := time.Unix(reminders[i].Timestamp, 0).String()
		buttonRow := []tb.InlineButton{tb.InlineButton{Unique: "delete_reminder", Text: messageText, Data: reminders[i].ID.Hex()}}
		buttonArray = append(buttonArray, buttonRow)
	}
	botInstance.Send(m.Sender, "Which reminder do you want to cancel?", &tb.ReplyMarkup{InlineKeyboard: buttonArray, ReplyKeyboardRemove: true})
}

func deleteReminderHandler(c *tb.Callback) {
	// TODO: Actually cancel reminder

	botInstance.Edit(c.Message, "Reminder cancelled!", &tb.ReplyMarkup{})
	botInstance.EditReplyMarkup(c.Message, &tb.ReplyMarkup{})
	fmt.Println(c.ID)
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
			go confirmReminderSet(wait, m.Sender)
			go forwardMessageAfterDelay(wait, m.Sender, waitingMessage)
		}

		delete(currentLimboUsers, m.Sender.ID)
	} else {
		currentLimboUsers[m.Sender.ID] = m
		botInstance.Send(m.Sender, "When should I remind you?")
	}
}
