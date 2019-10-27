package main

import (
	"net/http"
	"strings"
)

func fromTelegram(r *http.Request) bool {
	baseIP := "149.154.167"

	split := strings.Split(r.RemoteAddr, ".")
	if len(split) < 3 || strings.Join(split[0:3], ".") != baseIP {
		return false
	}
	return true
}
