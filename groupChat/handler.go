package main

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func init() {
	gc = make(groupChat, 0)
}

func hanleEvents(c *gin.Context) {
	var req request

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Fprintf(c.Writer, "Bind error: %+v\n", err)
		fmt.Printf("Bind error: %+v\n", err)
		c.JSON(500, err)
		return
	}

	if req.IsVerification() {
		c.Data(200, "text/plain", []byte(req.Challenge))
		return
	}

	fmt.Println(req)

	// ACK immediately
	c.Data(200, "text/plain", nil)

	// this is our bot's message, don't check
	if req.Event.ClientMsgID == "" {
		return
	}

	if req.Event.Type == "message" && req.Event.Text != "" {
		wrap := strings.Repeat("-", len(req.Event.Text)+2)
		fmt.Printf("%s\n|%s|\n%s\n", wrap, req.Event.Text, wrap)

		// 1. send to included users
		if invites, ok := gc[req.Event.Channel]; ok {
			for _, invite := range invites {
				postMessageToUser(req.Event.Text, invite.Channel, req.Event.User)
			}
			return
		}

		// 2. or reply to group the user is included to
		for groupChannel, invites := range gc {
			for _, invite := range invites {
				if invite.Channel == req.Event.Channel {
					// newChannel, _ := getUserChannel(invite.AddedBy)
					// if newChannel != "" {
					postMessageToGroup(req.Event.Text, groupChannel, req.Event.User)
					break // next group
					// }
				}
			}
		}
	}

}
