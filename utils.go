package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

func checkWebhook() bool {
	destURL := rootTelegramMethodURL + "/getWebhookInfo"
	res, err := http.Get(destURL)

	if err != nil {
		log.Println("Could not check webhook!")
	}
	bytes, _ := ioutil.ReadAll(res.Body)
	setURL := gjson.Get(string(bytes), "result.url").String()

	if setURL == "" {
		return false
	}
	return true
}

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
