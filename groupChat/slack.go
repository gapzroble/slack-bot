package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
}

func postMessageToUser(message, channel, sender string) error {
	return postMessageToSlack(message, channel, sender, botToken)
}

func postMessageToGroup(message, channel, sender string) error {
	return postMessageToSlack(message, channel, sender, slackToken)
}

func postMessageToSlack(message, channel, sender, token string) error {
	msg := gin.H{
		"text":    message,
		"channel": channel,
		"as_user": false,
	}
	if teamUser, ok := users[sender]; ok {
		msg["username"] = teamUser.Name
		msg["icon_url"] = teamUser.Image
	}

	r, _ := json.Marshal(msg)
	url := "https://slack.com/api/chat.postMessage"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(r))
	req.Header.Set("Content-type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error response:", err)
	}

	if resp != nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("body:", string(body))
	}

	return err
}

func getUserChannel(user string) (string, error) {
	m := gin.H{
		"token": slackToken,
		"user":  user,
	}
	r, _ := json.Marshal(m)
	url := "https://slack.com/api/im.open"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(r))
	req.Header.Set("Content-type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+botToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error response:", err)
		return "", err
	}

	var chat im
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("body:", string(body))

	if err := json.Unmarshal(body, &chat); err != nil {
		return "", err
	}

	return chat.Channel.ID, nil
}
