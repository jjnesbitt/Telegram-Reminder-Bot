package main

import (
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Wait represents a unit of time to wait
type Wait struct {
	units    string
	quantity int
	seconds  int
	duration time.Duration
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

var numKeyboardLayout = [][]string{
	{"7", "8", "9"},
	{"4", "5", "6"},
	{"1", "2", "3"},
	{"0", "=>"},
}
var replyBtn = tb.ReplyButton{Text: "🌕 Button #1"}
var numKeyboard = [][]tb.ReplyButton{
	[]tb.ReplyButton{replyBtn},
	// ...
}
