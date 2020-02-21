package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Load all users into memory
func Load(slackToken string) (Team, error) {
	var cursor string
	endpoint := "https://slack.com/api/users.list"
	data := url.Values{}
	data.Add("limit", "100")

	team := make(Team, 0)

	for {
		if cursor != "" {
			data.Set("cursor", cursor)
		}
		req, _ := http.NewRequest("POST", endpoint, bytes.NewBufferString(data.Encode()))
		req.Header.Set("Content-type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "Bearer "+slackToken)

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			fmt.Println("Error response:", err)
			return team, err
		}

		var result userlist

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("Error unmarshall:", err)
			return team, err
		}

		if result.ResponseMetadata.Next == "" {
			break
		}

		cursor = result.ResponseMetadata.Next

		for _, mem := range result.Members {
			if !mem.Deleted {
				newUser := user{
					Name:  mem.Profile.RealName,
					Image: mem.Profile.Image,
				}
				if mem.Profile.Name != "" {
					newUser.Name = mem.Profile.Name
				}
				team[mem.ID] = newUser
			}
		}
	}

	return team, nil
}
