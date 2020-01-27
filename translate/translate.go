package main

import (
	"encoding/json"
	"errors"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var translator = "https://translate.googleapis.com/translate_a/single?client=gtx&sl=auto&hl=sv&tl=en&dt=t&q="

func translate(message string) (trans string, e error) {
	replacements := extract(&message)
	defer replace(&trans, replacements)

	url := translator + url.QueryEscape(message)
	response, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to translate message from api: %s, Url: %s", err.Error(), url)
		return "", err
	}

	if response.StatusCode > 299 {
		log.Printf("Expecting response status 2XX: got %d", response.StatusCode)
		return "", errors.New("Translate response not OK")
	}

	dat, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Failed to read response body: %s", err.Error())
		return "", err
	}

	body := parse(string(dat))
	if body == "" {
		body = string(dat)
	}

	return body, nil
}

func extract(message *string) map[string]string {
	rep := make(map[string]string, 0)

	enc := false // <@Userxxx>
	enclosed := strings.FieldsFunc(*message, func(r rune) bool {
		switch r {
		case '<':
			enc = true
			return false
		case '>':
			enc = false
			return false
		}
		return !enc
	})
	for _, encl := range enclosed {
		if strings.HasPrefix(encl, "<") && strings.HasSuffix(encl, ">") {
			hash := crc32(encl)
			*message = strings.ReplaceAll(*message, encl, hash)
			rep[hash] = encl
		}
	}

	emoji := false // :emoji:
	emojis := strings.FieldsFunc(*message, func(r rune) bool {
		switch r {
		case ':':
			emoji = !emoji
			return false
		case ' ', '\n':
			if emoji {
				emoji = false
			}
		}
		return !emoji
	})
	for _, emo := range emojis {
		if len(emo) > 1 && strings.HasPrefix(emo, ":") && strings.HasSuffix(emo, ":") {
			hash := crc32(emo)
			*message = strings.ReplaceAll(*message, emo, hash)
			rep[hash] = emo
		}
	}

	return rep
}

func replace(message *string, replacements map[string]string) {
	for key, value := range replacements {
		*message = strings.ReplaceAll(*message, key, value)
	}
}

func crc32(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return "{" + strconv.FormatUint(uint64(h.Sum32()), 10) + "}"
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

	return translation
}
