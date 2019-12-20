package main

import (
	"fmt"
	"net/url"
	"strings"
)

func command(body string) (string, bool) {
	u, err := url.Parse("/?" + body)
	if err != nil {
		return "", false
	}

	q := u.Query()
	user := q.Get("user_id")
	if user == "" {
		return "", false
	}
	text := q.Get("text")
	what := strings.ToLower(strings.Trim(text, " "))
	switch what {
	case "":
		return ":point_right: " + status(user), true
	case "yes", "subscribe", "1", "true", "y":
		slackUsers[user] = true
		return status(user), true
	case "no", "unsubscribe", "0", "false", "n":
		delete(slackUsers, user)
		return status(user), true
	case "nerdstats":
		return fmt.Sprintf("```slackUsers: %+v\nnotInChannel: %+v```", slackUsers, notInChannel), true
	}

	message, err := translate(text)
	if err != nil {
		return err.Error(), true
	}

	if strings.ToLower(message) == strings.ToLower(text) {
		return text, true
	}

	return message, true
}

func status(user string) string {
	if _, ok := slackUsers[user]; ok {
		return "Subscribed :heavy_check_mark:"
	}
	return "~Unsubscribed~"
}
