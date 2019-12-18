package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/tiqqe/go-logger"
)

var slackToken string

func init() {
	slackToken = os.Getenv("SLACK_API_TOKEN")
}

func postMessageToSlack(message, channel, sender, user, ts string) error {
	msg := map[string]interface{}{
		"text":    message,
		"channel": channel,
		// "as_user":  true,
		// "username": sender,
		"user": user,
		// "thread_ts": ts, // TODO: identify if reply
	}
	logger.Info(&logger.LogEntry{
		Message: "Sending message",
		Keys: map[string]interface{}{
			"Message": msg,
		},
	})
	r, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	url := "https://slack.com/api/chat.postEphemeral"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(r))
	req.Header.Set("Content-type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+slackToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	logger.Info(&logger.LogEntry{
		Message: "Got response from slack",
		Keys: map[string]interface{}{
			"Response": string(body),
		},
	})

	res, err := newResponse(body)
	if err != nil {
		return err
	}
	if res.Error == "user_not_in_channel" {
		return errors.New("user_not_in_channel")
	}

	return nil
}

func getUserChannel(user string) (string, error) {
	msg := map[string]interface{}{
		"token": slackToken,
		"user":  user,
	}
	r, _ := json.Marshal(msg)
	url := "https://slack.com/api/im.open"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(r))
	req.Header.Set("Content-type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+slackToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	var chat im
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &chat); err != nil {
		return "", err
	}

	return chat.Channel.ID, nil
}
