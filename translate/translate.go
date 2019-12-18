package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/tiqqe/go-logger"
)

var translator = "https://translate.googleapis.com/translate_a/single?client=gtx&sl=sv&tl=en&dt=t&q="

func translate(message string) (string, error) {
	url := translator + url.QueryEscape(message)
	response, err := http.Get(url)
	if err != nil {
		logger.Error(&logger.LogEntry{
			Message:      "Failed to translate message from api",
			ErrorMessage: err.Error(),
			Keys: map[string]interface{}{
				"Url": url,
			},
		})
		return "", err
	}

	if response.StatusCode > 299 {
		logger.Error(&logger.LogEntry{
			Message: "Expecting response status 2XX",
			Keys: map[string]interface{}{
				"Response": response,
			},
		})
		return "", errors.New("Translate response not OK")
	}

	dat, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error(&logger.LogEntry{
			Message:      "Fail to read response body",
			ErrorMessage: err.Error(),
		})
		return "", err
	}

	body := parse(string(dat))
	if body == "" {
		body = string(dat)
	}

	return body, nil
}

func parse(response string) string {
	var trans interface{}
	if err := json.Unmarshal([]byte(response), &trans); err != nil {
		return ""
	}

	root, ok := trans.([]interface{})
	if !ok {
		return ""
	}

	translations, ok := root[0].([]interface{})
	if !ok {
		return ""
	}

	translation := ""
	for _, tr := range translations {
		actual, ok := tr.([]interface{})
		if !ok {
			continue
		}

		str, ok := actual[0].(string)
		if ok {
			translation += str
		}
	}

	return strings.ReplaceAll(translation, "<@ U", "<@U")
}
