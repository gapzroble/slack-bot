package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func getSenderPic(sender string) string {
	url := fmt.Sprintf("https://slack.com/api/users.info?token=%s&user=%s", slackToken, sender)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error getting user.info, %s", err.Error())
		return ""
	}

	var info UserInfo
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading body, %s", err.Error())
		return ""
	}
	// log.Printf("body: %s", body)

	if err := json.Unmarshal(body, &info); err != nil {
		log.Printf("Error unmarshalling body, %s", err.Error())
		return ""
	}

	if !info.OK {
		log.Printf("Info not OK, %#v", info)
		return ""
	}

	return info.User.Profile.Image48
}
