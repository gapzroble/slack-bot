package main

// UserInfo struct
type UserInfo struct {
	OK   bool `json:"ok"`
	User struct {
		Profile struct {
			Image48 string `json:"image_48"`
		} `json:"profile"`
	} `json:"user"`
}
