package slack

import (
	"encoding/json"
)

// Action struct
type Action struct {
	CallbackID string `json:"callback_id"`
	User       struct {
		ID string `json:"id"`
	} `json:"user"`
	Channel struct {
		ID string `json:"id"`
	} `json:"channel"`
	Message struct {
		Text string `json:"text"`
		User string `json:"user"`
		TS   string `json:"ts"`
	} `json:"message"`
	TriggerID string `json:"trigger_id"`
}

// NewAction func
func NewAction(s string) (*Action, error) {
	var action Action
	if err := json.Unmarshal([]byte(s), &action); err != nil {
		return nil, err
	}

	return &action, nil
}
