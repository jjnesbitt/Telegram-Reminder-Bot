package main

import (
	"log"
	"net/http"
	"net/url"
)

func setWebhook() {
	destURL := rootTelegramMethodURL + "/setWebhook"
	_, err := http.PostForm(destURL, url.Values{"url": {publicURL}})

	if err != nil {
		log.Fatal("Error setting webhook: " + err.Error())
	} else {
		log.Println("Webhook set successfully")
	}
}

func deleteWebhook() {
	destURL := rootTelegramMethodURL + "/setWebhook"
	_, err := http.PostForm(destURL, url.Values{"url": {""}})

	if err != nil {
		log.Fatal("Error unsetting webhook: " + err.Error())
	} else {
		log.Println("Webhook unset successfully")
	}
}
