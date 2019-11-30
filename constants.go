package main

import (
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Wait represents a unit of time to wait
type Wait struct {
	units           string
	quantity        int
	seconds         int64
	duration        time.Duration
	futureTimestamp int64
}

// StoredReminder stores a message
type StoredReminder struct {
	ChatID    int64    `bson:"chat_id"`
	MessageID int      `bson:"message_id"`
	User      *tb.User `bson:"user"`
	Time      int64    `bson:"timestamp"`
}

var timeUnits = []string{"second", "minute", "hour", "day", "week", "month"}

// Maps unit names to seconds
var unitMap = map[string]int{
	"second": 1,
	"minute": 60,
	"hour":   3600,
	"day":    86400,
	"week":   604800,
	"month":  2419200000,
}

var keyboardTimeUnitButtons = []tb.ReplyButton{
	tb.ReplyButton{Text: "sec"},
	tb.ReplyButton{Text: "min"},
	tb.ReplyButton{Text: "hr"},
	tb.ReplyButton{Text: "day"},
	tb.ReplyButton{Text: "week"},
	tb.ReplyButton{Text: "month"},
}

var numKeyboardLayout = [][]tb.ReplyButton{
	{tb.ReplyButton{Text: "7"}, tb.ReplyButton{Text: "8"}, tb.ReplyButton{Text: "9"}},
	{tb.ReplyButton{Text: "4"}, tb.ReplyButton{Text: "5"}, tb.ReplyButton{Text: "6"}},
	{tb.ReplyButton{Text: "1"}, tb.ReplyButton{Text: "2"}, tb.ReplyButton{Text: "3"}},
	{tb.ReplyButton{Text: "0"}},
	// append([]tb.ReplyButton{tb.ReplyButton{Text: "0"}}, keyboardTimeUnitButtons...),
	keyboardTimeUnitButtons,
}

var numKeyboardMarkup = tb.ReplyMarkup{ReplyKeyboard: numKeyboardLayout}
