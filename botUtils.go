package main

import (
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	tb "gopkg.in/tucnak/telebot.v2"
)

// botInstance is the global instance of tb.Bot that all functions use
var botInstance *tb.Bot

func getBotPreferences() tb.Settings {
	listenPort := ":" + envListenPort

	var poller tb.Poller = &tb.Webhook{
		Listen:   listenPort,
		Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
	}

	if envMode != "production" {
		deleteWebhook()
		poller = &tb.LongPoller{Timeout: 1 * time.Second}
	} else {
		if !webhookIsSet() {
			setWebhook()
		}
	}

	pref := tb.Settings{
		Token:  envBotToken,
		Poller: poller,
	}

	return pref
}

func confirmReminderSet(wait Reminder, recipient tb.Recipient) {
	stringQuantity := strconv.Itoa(wait.quantity)
	string := "Reminder set for " + stringQuantity + " " + wait.units + "s!"
	botInstance.Send(recipient, string)
}

func forwardMessage(recipient *tb.User, message *tb.Message) {
	botInstance.Forward(recipient, message)
}

func forwardStoredMessageAfterDelay(id primitive.ObjectID, duration time.Duration) {
	time.Sleep(duration)

	rem, err := getStoredReminderFromID(id)
	if err != nil {
		log.Println("Message unable to be retrieved or no longer in database")
		return
	}

	message := messageFromStoredReminder(rem)

	go forwardMessage(rem.User, &message)
	go removeMessageFromDB(id)
}

func forwardMessageAfterDelay(wait Reminder, recipient *tb.User, message *tb.Message) {
	id := storeMessageIntoDB(message, recipient, wait.timestamp)
	forwardStoredMessageAfterDelay(id, wait.duration)
}

func getWaitTime(payload string) (Reminder, error) {
	temp := strings.Join(timeUnits, "|")
	waitExpr := regexp.MustCompile(`(\d+) (` + temp + `)s?`)

	matches := waitExpr.FindStringSubmatch(payload)

	if matches == nil {
		return Reminder{}, errors.New("No matches found")
	}

	quant, _ := strconv.Atoi(matches[1])
	units := matches[2]

	// TODO: Fix how duration is being generated
	seconds := int64(quant * unitMap[units])
	duration := time.Duration(seconds) * time.Second
	timestamp := time.Now().Add(duration)
	return Reminder{units: units, quantity: quant, duration: duration, timestamp: timestamp.Unix()}, nil
}
