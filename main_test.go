package main

import (
	"testing"
)

func TestGetWaitTime(t *testing.T) {
	match, err := getWaitTime("20 hours")
	if err {
		t.Error("Valid Wait match not found")
	} else if match.quantity != 20 || match.units != "hour" {
		t.Error("Incorrect Wait values returned")
	}

	match, err = getWaitTime("15 not units")
	if !err {
		t.Error("Valid Wait matched from invalid string")
	}
}
