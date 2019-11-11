package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

func getBotPreferences() tb.Settings {
	botToken := os.Getenv("BOT_TOKEN")
	listenPort := ":" + os.Getenv("LISTEN_PORT")
	publicURL := os.Getenv("PUBLIC_URL")

	webhook := &tb.Webhook{
		Listen:   listenPort,
		Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
	}

	pref := tb.Settings{
		Token:  botToken,
		Poller: webhook,
	}

	return pref
}

func confirmReminderSet(wait *Wait, b *tb.Bot, sender *tb.User) {
	stringQuantity := strconv.Itoa(wait.quantity)
	string := "Reminder set for " + stringQuantity + " " + wait.units + "s!"
	b.Send(sender, string)
}

func sendMessage(m *tb.Message) {

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

	return Wait{units: units, quantity: quant, seconds: seconds}, false
}
