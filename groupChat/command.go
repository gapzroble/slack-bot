package main

import (
	"net/url"
)

func doCommand(body string) (string, bool) {
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

	return text, true
}
