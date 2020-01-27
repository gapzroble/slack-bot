package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var slackToken string

func init() {
	slackToken = os.Getenv("SLACK_API_TOKEN")
}

func quote(s string) string {
	return ">" + strings.ReplaceAll(s, "\n", "\n>")
}

func postMessageToSlack(message, channel, sender, user, ts string) error {
	msg := map[string]interface{}{
		"text":    quote(message), // quote it
		"channel": channel,
		"user":    user,
	}
	if ts != "" {
		msg["thread_ts"] = ts
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
	log.Printf("Got response from slack: %s", body)

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

func getUserChannel(user string) (string, error) {
	msg := map[string]interface{}{
		"token": slackToken,
		"users": user,
	}
	r, _ := json.Marshal(msg)
	url := "https://slack.com/api/conversations.open"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(r))
	req.Header.Set("Content-type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+slackToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error getting user channel: %s", err.Error())
		return "", err
	}

	var chat im
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &chat); err != nil {
		log.Printf("Error unmarshalling user channel response: %s", err.Error())
		return "", err
	}

	return chat.Channel.ID, nil
}

func showModal(dat []byte) error {
	url := "https://slack.com/api/views.open"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dat))
	req.Header.Set("Content-type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+slackToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error opening view: %s", err.Error())
		return err
	}

	var chat im
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &chat); err != nil {
		log.Printf("Error unmarshalling view response: %s", err.Error())
		return err
	}
	log.Printf("Got response from slack: %s", body)

	return nil
}
