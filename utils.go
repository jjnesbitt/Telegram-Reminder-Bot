package main

import (
	"regexp"
	"strconv"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

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
