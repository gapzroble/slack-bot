package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/tiqqe/go-logger"
)

var slackToken string
var slackUser string

func init() {
	slackToken = os.Getenv("SLACK_API_TOKEN")
	slackUser = os.Getenv("SLACK_USER")
}

func postMessageToSlack(message, channel, sender, user, ts string) error {
	msg := map[string]interface{}{
		"text":    message,
		"channel": channel,
		// "as_user":  true,
		// "username": sender,
		"user": user,
		// "thread_ts": ts,
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

	body, e := ioutil.ReadAll(resp.Body)
	if e == nil {
		defer resp.Body.Close()
		logger.Error(&logger.LogEntry{
			Message: "Got response from slack",
			Keys: map[string]interface{}{
				"Response": string(body),
			},
		})
	}

	return err
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
