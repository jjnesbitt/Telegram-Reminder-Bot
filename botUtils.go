package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	tb "gopkg.in/tucnak/telebot.v2"
)

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

func confirmReminderSet(wait *Wait, b *tb.Bot, recipient tb.Recipient) {
	stringQuantity := strconv.Itoa(wait.quantity)
	string := "Reminder set for " + stringQuantity + " " + wait.units + "s!"
	b.Send(recipient, string)
}

func forwardMessage(b *tb.Bot, recipient *tb.User, message *tb.Message) {
	b.Forward(recipient, message)
}

// TODO: Make all of thse store and pull from database each time instead of holding onto references
func forwardStoredMessageAfterDelay(id primitive.ObjectID, wait Wait, b *tb.Bot, recipient *tb.User, message *tb.Message) {
	time.Sleep(wait.duration)
	go forwardMessage(b, recipient, message)
	go removeMessageFromDB(id)
}

func forwardMessageAfterDelay(wait Wait, b *tb.Bot, recipient *tb.User, message *tb.Message) {
	id := storeMessageIntoDB(message, recipient, wait)
	forwardStoredMessageAfterDelay(id, wait, b, recipient, message)
}

func getWaitTime(payload string) (Wait, error) {
	temp := strings.Join(timeUnits, "|")
	waitExpr := regexp.MustCompile(`(\d+) (` + temp + `)s?`)

	matches := waitExpr.FindStringSubmatch(payload)

	if matches == nil {
		return Wait{}, errors.New("No matches found")
	}

	quant, _ := strconv.Atoi(matches[1])
	units := matches[2]
	seconds := int64(quant * unitMap[units])
	duration := time.Duration(seconds * int64(time.Second))
	futureTimestamp := time.Now().Unix() + int64(duration.Seconds())

	return Wait{units: units, quantity: quant, seconds: seconds, duration: duration, futureTimestamp: futureTimestamp}, nil
}
