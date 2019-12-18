package main

import "encoding/json"

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

type with struct {
	Token       string `form:"token"`
	TeamID      string `form:"team_id"`
	ChannelID   string `form:"channel_id"`
	ChannelName string `form:"channel_name"`
	UserID      string `form:"user_id"`
	UserName    string `form:"user_name"`
	Command     string `form:"command"`
	Text        string `form:"text"`
	ResponseURL string `form:"response_url"`
	TriggerID   string `form:"trigger_id"`
}

type channel struct {
	ID string `json:"id"`
}
type im struct {
	Channel channel `json:"channel"`
}
