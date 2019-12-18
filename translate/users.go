package main

import (
	"os"
	"strings"
)

var notInChannel map[string]bool
var slackUsers map[string]bool

func init() {
	notInChannel = make(map[string]bool)
	slackUsers = make(map[string]bool)
	for _, user := range strings.Split(os.Getenv("SLACK_USER"), ",") {
		slackUsers[user] = true
	}
}

func userNotInChannel(user, channel string) bool {
	val, ok := notInChannel[user+":"+channel]
	if !ok {
		return false
	}

	return val
}
