package main

import (
	"testing"
)

func TestGetWaitTime(t *testing.T) {
	match, err := getWaitTime("20 hours")
	if err != nil {
		t.Error("Valid Reminder match not found")
	} else if match.quantity != 20 || match.units != "hour" {
		t.Error("Incorrect Reminder values returned")
	}

	match, err = getWaitTime("15 not units")
	if err == nil {
		t.Error("Valid Reminder matched from invalid string")
	}
}
