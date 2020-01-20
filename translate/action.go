package main

import (
	"encoding/json"
	"fmt"
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

	threadChan := make(chan string)
	go getMainThread(act.Message.TS, act.Channel.ID, threadChan)

	trans, err := translate(act.Message.Text)
	if err != nil {
		logger.ErrorStringf("Error translation, %s", err.Error())
		return "", false
	}

	threadTs := <-threadChan

	msg := fmt.Sprintf("<@%s>: %s", act.Message.User, trans)
	if err := postMessageToSlack(msg, act.Channel.ID, "bot", act.User.ID, threadTs); err != nil {
		logger.ErrorStringf("Error posting translation to user, %s", err.Error())
		return trans, true
	}

	return trans, true
}
