package main

type groupChat map[string][]chatInvite

type chatInvite struct {
	User    string //userId
	Channel string
	AddedBy string //userId
}

// addGroupChat and return user's channel id
func addToGc(groupChannel, user, by string) string {
	if _, ok := gc[groupChannel]; !ok {
		gc[groupChannel] = make([]chatInvite, 0)
	}

	for _, invite := range gc[groupChannel] {
		if invite.User == user {
			return invite.Channel
		}
	}

	userChannel, err := getUserChannel(user)
	if err != nil {
		return ""
	}

	newInvite := chatInvite{
		User:    user,
		Channel: userChannel,
		AddedBy: by,
	}

	gc[groupChannel] = append(gc[groupChannel], newInvite)

	return userChannel
}

func (g *groupChat) removeFromGc(groupChannel, user string) string {
	invites, ok := gc[groupChannel]
	if !ok {
		return ""
	}

	for i, invite := range invites {
		if invite.User == user {
			gc[groupChannel] = append(invites[0:i], invites[i+1:]...)
			return invite.Channel
		}
	}

	return ""

}
