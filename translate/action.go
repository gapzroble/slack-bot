package main

import "net/url"

import "encoding/json"

import "github.com/tiqqe/go-logger"

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
	} `json:"message"`
}

func newAction(s string) (*messageAction, error) {
	var act messageAction
	if err := json.Unmarshal([]byte(s), &act); err != nil {
		return nil, err
	}

	return &act, nil
}

func action(body string) (string, bool) {
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

	trans, err := translate(act.Message.Text)
	if err != nil {
		logger.ErrorStringf("Error translation, %s", err.Error())
		return "", false
	}

	return trans, true
}
