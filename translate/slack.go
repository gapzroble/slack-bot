package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
		"user":      user,
		"thread_ts": ts,
	}
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

func getPermalink(ts, channel string) (*permalink, error) {
	url := fmt.Sprintf("https://slack.com/api/chat.getPermalink?token=%s&channel=%s&message_ts=%s", slackToken, channel, ts)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var perm permalink
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &perm); err != nil {
		return nil, err
	}

	return &perm, nil
}
