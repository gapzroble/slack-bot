package users

type userlist struct {
	Members          []member `json:"members"`
	ResponseMetadata meta     `json:"response_metadata"`
}

type member struct {
	ID      string  `json:"id"`
	Profile profile `json:"profile"`
	Deleted bool    `json:"deleted"`
}

type profile struct {
	Name     string `json:"display_name"`
	RealName string `json:"real_name"`
	Image    string `json:"image_48"`
}

type meta struct {
	Next string `json:"next_cursor"`
}

type user struct {
	Name  string
	Image string
}

// Team type for users in memory
type Team map[string]user
