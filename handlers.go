package main

import (
	"log"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	tb "gopkg.in/tucnak/telebot.v2"
)

func setHandlers() {
	botInstance.Handle("/remindme", remindMeHandler)
	botInstance.Handle("/cancel", cancelHandler)
	botInstance.Handle("/list", listHandler)

	botInstance.Handle(tb.OnText, onTextHandler)
	botInstance.Handle(tb.OnCallback, callbackHandler)
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
	if privateMessageHelper(m, false) {
		return
	}

	reminders, done := reminderHelper(m)

	if !done {
		var buttonArray [][]tb.InlineButton

		buttonArray = append(buttonArray, []tb.InlineButton{tb.InlineButton{Unique: CallbackCancelAll, Text: "Cancel All"}})

		for i := range reminders {
			messageText := time.Unix(reminders[i].Timestamp, 0).String()
			buttonRow := []tb.InlineButton{tb.InlineButton{Unique: CallbackCancelReminder, Text: messageText, Data: reminders[i].ID.Hex()}}
			buttonArray = append(buttonArray, buttonRow)
		}
		buttonArray = append(buttonArray, []tb.InlineButton{tb.InlineButton{Unique: CallbackAbortCancel, Text: "Abort"}})
		botInstance.Send(m.Sender, "Which reminder do you want to cancel?", &tb.ReplyMarkup{InlineKeyboard: buttonArray, ReplyKeyboardRemove: true})
	}
}

func listHandler(m *tb.Message) {
	if privateMessageHelper(m, false) {
		return
	}

	reminders, done := reminderHelper(m)

	if !done {
		textArray := []string{"You have " + strconv.Itoa(len(reminders)) + " active reminders:"}
		for i := range reminders {
			textArray = append(textArray, time.Unix(reminders[i].Timestamp, 0).String())
		}

		text := strings.Join(textArray, "\n")
		botInstance.Send(m.Sender, text)
	}
}

// Handles direct forwarded requests
func onTextHandler(m *tb.Message) {
	if privateMessageHelper(m, true) {
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

func callbackHandler(c *tb.Callback) {
	data := strings.TrimSpace(c.Data)

	identifiers := strings.Split(data, "|")
	unique := identifiers[0]

	switch unique {
	case CallbackCancelAll:
		removeAllUserMessages(c.Sender)
		botInstance.Edit(c.Message, "All reminders cancelled!", &tb.ReplyMarkup{})

	case CallbackCancelReminder:
		idStr := identifiers[1]
		id, _ := primitive.ObjectIDFromHex(idStr)

		removeMessageFromDB(id)
		botInstance.Edit(c.Message, "Reminder cancelled!", &tb.ReplyMarkup{})

	case CallbackAbortCancel:
		botInstance.Edit(c.Message, "No reminders cancelled", &tb.ReplyMarkup{})
	}
}
