package main

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/rroble/slack-bot/lib/google"
	"github.com/rroble/slack-bot/lib/slack"
)

func doAction(body string) (string, bool) {
	u, err := url.Parse("/?" + body)
	if err != nil {
		log.Printf("Error parsing body: %s", err.Error())
		return "", false
	}

	payload := u.Query().Get("payload")
	if payload == "" {
		log.Print("Payload is empty")
		return "", false
	}

	act, err := slack.NewAction(payload)
	if err != nil {
		log.Printf("Error new action: %s", err.Error())
		return "", false
	}
	log.Printf("Got action: %#v", act)

	trans, err := google.Translate(act.Message.Text)
	if err != nil {
		log.Printf("Error translation: %s", err.Error())
		return "", false
	}

	mod := slack.NewModal(act.TriggerID, act.Message.Text, trans)
	dat, err := json.Marshal(mod)
	if err != nil {
		log.Printf("Error marshalling modal: %s", err.Error())
		return trans, true
	}

	if err := slack.ShowModal(dat); err != nil {
		log.Printf("Error showing modal: %s", err.Error())
	}

	return trans, true
}
