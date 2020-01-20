package main

import (
	"encoding/json"
	"net/url"

	"github.com/tiqqe/go-logger"
)

type messageAction struct {
	CallbackID string `json:"callback_id"`
	User       struct {
		ID string `json:"id"`
	} `json:"user"`
	Channel struct {
		ID string `json:"id"`
	} `json:"channel"`
	Message struct {
		Text string `json:"text"`
		User string `json:"user"`
		TS   string `json:"ts"`
	} `json:"message"`
	TriggerID string `json:"trigger_id"`
}

func newAction(s string) (*messageAction, error) {
	var act messageAction
	if err := json.Unmarshal([]byte(s), &act); err != nil {
		return nil, err
	}

	return &act, nil
}

func doAction(body string) (string, bool) {
	u, err := url.Parse("/?" + body)
	if err != nil {
		logger.ErrorStringf("Error parsing body, %s", err.Error())
		return "", false
	}

	payload := u.Query().Get("payload")
	if payload == "" {
		logger.ErrorString("Payload is empty")
		return "", false
	}

	act, err := newAction(payload)
	if err != nil {
		logger.ErrorStringf("Error new action, %s", err.Error())
		return "", false
	}
	logger.Info(&logger.LogEntry{
		Message: "Got action",
		Keys: map[string]interface{}{
			"Action": act,
		},
	})

	trans, err := translate(act.Message.Text)
	if err != nil {
		logger.ErrorStringf("Error translation, %s", err.Error())
		return "", false
	}

	mod := newModal(act.TriggerID, act.Message.Text, trans)
	dat, err := json.Marshal(mod)
	if err != nil {
		logger.ErrorStringf("Error marshalling modal, %s", err.Error())
		return trans, true
	}

	if err := showModal(dat); err != nil {
		logger.ErrorStringf("Error showing modal, %s", err.Error())
	}

	return trans, true
}
