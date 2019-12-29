package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	tb "gopkg.in/tucnak/telebot.v2"
)

// Unique identifiers for various inline callbacks
const (
	// Selected reminder to cancel
	CallbackCancelReminder = "cancel_reminder"

	// To abort the /cancel command
	CallbackAbortCancel = "abort_cancel"
)

// Reminder represents a unit of time to wait
type Reminder struct {
	units     string
	quantity  int
	duration  time.Duration
	timestamp int64
}

// StoredReminder stores a message
type StoredReminder struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ChatID    int64              `bson:"chat_id"`
	MessageID int                `bson:"message_id"`
	User      *tb.User           `bson:"user"`
	Timestamp int64              `bson:"timestamp"`
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
