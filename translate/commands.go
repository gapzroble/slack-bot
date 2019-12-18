package main

import (
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
		return "*" + status(user), true
	case "yes":
		slackUsers[user] = true
		return status(user), true
	case "no":
		delete(slackUsers, user)
		return status(user), true
	}

	message, err := translate(text)
	if err != nil {
		return err.Error(), true
	}

	return message, true
}

func status(user string) string {
	if _, ok := slackUsers[user]; ok {
		return "Subscribed"
	}
	return "Unsubscribed"
}
