package main

import (
	"regexp"
	"strconv"
	"strings"
	"time"

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
		if !webhookSet() {
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

func sendMessage(m *tb.Message) {

}

func forwardMessageAfterDelay(delay time.Duration, b *tb.Bot, recipient *tb.User, message *tb.Message) {
	time.Sleep(delay)
	b.Forward(recipient, message)
}

func getWaitTime(payload string) (Wait, bool) {
	temp := strings.Join(timeUnits, "|")
	waitExpr := regexp.MustCompile(`(\d+) (` + temp + `)s?`)

	matches := waitExpr.FindStringSubmatch(payload)

	if matches == nil {
		return Wait{}, true
	}

	quant, _ := strconv.Atoi(matches[1])
	units := matches[2]
	seconds := quant * unitMap[units]
	duration := time.Duration(int64(seconds * int(time.Second)))

	return Wait{units: units, quantity: quant, seconds: seconds, duration: duration}, false
}
