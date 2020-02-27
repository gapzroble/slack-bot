package slack

import (
	"encoding/json"
	"strings"
)

// Reaction struct
type Reaction struct {
	Name  string
	Count int
	Users []string
}

// Event struct
type Event struct {
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
	Reactions   []Reaction `json:",omitempty"`
}

// Request struct
type Request struct {
	Token       string
	TeamID      string `json:"team_id"`
	APIAppID    string `json:"api_app_id"`
	Event       Event
	Type        string
	EventID     string   `json:"event_id"`
	EventTime   float64  `json:"event_time"`
	AuthedUsers []string `json:"authed_users"`
	Challenge   string   `json:",omitempty"`
}

// IsVerification func
func (r Request) IsVerification() bool {
	return r.Type == "url_verification" && r.Challenge != ""
}

// String func
func (r Request) String() string {
	s, _ := json.MarshalIndent(r, "", "\t")
	return string(s)
}

// Channel struct
type Channel struct {
	ID string `json:"id"`
}

// IM struct
type IM struct {
	Channel Channel `json:"channel"`
}

// Response struct
type Response struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

// NewResponse from byte array
func NewResponse(data []byte) (*Response, error) {
	var res Response
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Permalink struct
type Permalink struct {
	OK        bool   `json:"ok"`
	Channel   string `json:"channel"`
	Permalink string `json:"permalink"`
}

// Modal struct
type Modal struct {
	Token     string `json:"token"`
	TriggerID string `json:"trigger_id"`
	View      string `json:"view"`
}

// NewModal func
func NewModal(triggerID, message, translation string) Modal {
	r := strings.NewReplacer("{message}", message, "{translation}", translation)
	return Modal{
		TriggerID: triggerID,
		View:      r.Replace(modalView),
	}
}

// UserInfo struct
type UserInfo struct {
	OK   bool `json:"ok"`
	User struct {
		Profile struct {
			Image48 string `json:"image_48"`
		} `json:"profile"`
	} `json:"user"`
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
                "type": "mrkdwn",
                "text": "{message}"
            }
        },
        {
            "type": "divider"
        },
        {
            "type": "section",
            "text": {
                "type": "mrkdwn",
                "text": "{translation}"
            }
        }
    ]
}
`
