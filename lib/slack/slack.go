package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	slackToken          string
	postMessageURL      = "https://slack.com/api/chat.postEphemeral"
	getPermalinkURL     = "https://slack.com/api/chat.getPermalink?token={token}&channel={channel}&message_ts={message_ts}"
	conversationOpenURL = "https://slack.com/api/conversations.open"
	viewOpenURL         = "https://slack.com/api/views.open"
	userInfoURL         = "https://slack.com/api/users.info?token={token}&user={user}"
)

func init() {
	slackToken = os.Getenv("SLACK_API_TOKEN")
}

// PostMessage func
func PostMessage(message, channel, sender, user string, tsChan, senderPicChan <-chan string) error {
	msg := map[string]interface{}{
		"text":    message,
		"channel": channel,
		"user":    user,
	}
	senderPic := <-senderPicChan
	if senderPic != "" {
		msg["as_user"] = false
		msg["icon_url"] = senderPic
	}

	ts := <-tsChan
	if ts != "" {
		msg["thread_ts"] = ts
	}

	r, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", postMessageURL, bytes.NewBuffer(r))
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

	res, err := NewResponse(body)
	if err != nil {
		return err
	}

	if res.Error == "user_not_in_channel" {
		return errors.New("user_not_in_channel")
	}

	return nil
}

// GetPermalink func
func GetPermalink(ts, channel string) (*Permalink, error) {
	r := strings.NewReplacer("{token}", slackToken, "{channel}", channel, "{message_ts}", ts)
	resp, err := http.Get(r.Replace(getPermalinkURL))
	if err != nil {
		return nil, err
	}

	var perm Permalink
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

// GetUserChannel func
func GetUserChannel(user string) (string, error) {
	msg := map[string]interface{}{
		"token": slackToken,
		"users": user,
	}
	r, _ := json.Marshal(msg)
	req, err := http.NewRequest("POST", conversationOpenURL, bytes.NewBuffer(r))
	req.Header.Set("Content-type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+slackToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error getting user channel: %s", err.Error())
		return "", err
	}

	var chat IM
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &chat); err != nil {
		log.Printf("Error unmarshalling user channel response: %s", err.Error())
		return "", err
	}

	return chat.Channel.ID, nil
}

// ShowModal func
func ShowModal(dat []byte) error {
	req, err := http.NewRequest("POST", viewOpenURL, bytes.NewBuffer(dat))
	req.Header.Set("Content-type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+slackToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error opening view: %s", err.Error())
		return err
	}

	var chat IM
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &chat); err != nil {
		log.Printf("Error unmarshalling view response: %s", err.Error())
		return err
	}
	log.Printf("Got response from slack: %s", body)

	return nil
}

// GetMainThread func
func GetMainThread(ts, channel string) string {
	perm, err := GetPermalink(ts, channel)
	if err != nil {
		return ""
	}

	if !perm.OK {
		return ""
	}

	u, err := url.Parse(perm.Permalink)
	if err != nil {
		return ""
	}

	tsparam := u.Query().Get("thread_ts")
	if tsparam == "" || tsparam != ts {
		return tsparam
	}

	return ""
}

// GetSenderPic func
func GetSenderPic(sender string) string {
	r := strings.NewReplacer("{token}", slackToken, "{user}", sender)
	resp, err := http.Get(r.Replace(userInfoURL))
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
