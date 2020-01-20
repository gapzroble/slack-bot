package main

import "encoding/json"

import "strings"

type reaction struct {
	Name  string
	Count int
	Users []string
}

type event struct {
	ClientMsgID string `json:"client_msg_id"`
	Type        string
	SubtType    string `json:"subtype,omitempty"`
	Text        string
	User        string
	TS          string `json:"ts"`
	Team        string
	Channel     string
	EventTS     string     `json:"event_ts"`
	ChannelType string     `json:"channel_type"`
	Hidden      bool       `json:",omitempty"`
	DeleteTS    string     `json:"deleted_ts,omitempty"`
	IsStarred   bool       `json:"is_starred,omitempty"`
	PinnedTo    []string   `json:"pinned_to,omitempty"`
	Reactions   []reaction `json:",omitempty"`
}

type request struct {
	Token       string
	TeamID      string `json:"team_id"`
	APIAppID    string `json:"api_app_id"`
	Event       event
	Type        string
	EventID     string   `json:"event_id"`
	EventTime   float64  `json:"event_time"`
	AuthedUsers []string `json:"authed_users"`
	Challenge   string   `json:",omitempty"`
}

func (r request) IsVerification() bool {
	return r.Type == "url_verification" && r.Challenge != ""
}

func (r request) String() string {
	s, _ := json.MarshalIndent(r, "", "\t")
	return string(s)
}

type channel struct {
	ID string `json:"id"`
}
type im struct {
	Channel channel `json:"channel"`
}

type response struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

func newResponse(data []byte) (*response, error) {
	var res response
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

type permalink struct {
	OK        bool   `json:"ok"`
	Channel   string `json:"channel"`
	Permalink string `json:"permalink"`
}

type modal struct {
	Token     string `json:"token"`
	TriggerID string `json:"trigger_id"`
	View      string `json:"view"`
}

func newModal(triggerID, message, translation string) modal {
	view := strings.ReplaceAll(modalView, "--message--", message)
	view = strings.ReplaceAll(view, "--translation--", translation)
	return modal{
		TriggerID: triggerID,
		View:      view,
	}
}

var modalView = `
{
    "type": "modal",
    "title": {
        "type": "plain_text",
        "text": "Translate Bot",
        "emoji": true
    },
    "blocks": [
        {
            "type": "section",
            "text": {
                "type": "plain_text",
                "text": "--message--",
                "emoji": true
            }
        },
        {
            "type": "divider"
        },
        {
            "type": "section",
            "text": {
                "type": "plain_text",
                "text": "--translation--",
                "emoji": true
            }
        }
    ]
}
`
