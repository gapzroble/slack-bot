package main

var notInChannel map[string]bool

func init() {
	if len(notInChannel) == 0 {
		notInChannel = make(map[string]bool)
	}
}

func userNotInChannel(user, channel string) bool {
	val, ok := notInChannel[user+channel]
	if !ok {
		return false
	}

	return val
}
