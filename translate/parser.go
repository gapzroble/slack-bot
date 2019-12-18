package main

import (
	"encoding/json"
	"strings"
)

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
