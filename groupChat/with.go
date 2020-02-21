package main

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func hanleWith(c *gin.Context) {
	var req with

	if err := c.ShouldBindWith(&req, binding.FormPost); err != nil {
		fmt.Println("Bind error:", err)
		c.JSON(500, err)
		return
	}

	if !req.Allowed() {
		c.JSON(200, gin.H{
			"text": "Only allowed in private group",
		})
		return
	}

	fmt.Printf("%+v\n", req)

	users := getUsers(req.Text)
	if len(users) == 0 {
		c.Data(200, "text/plain", nil)
		fmt.Println("no users found", req.Text)
		return
	}

	// return asap
	c.JSON(200, gin.H{
		"text":          "",
		"response_type": "in_channel",
	})

	fmt.Println("found users", users)

	for _, user := range users {
		userChannel := addToGc(req.ChannelID, user, req.UserID)
		if userChannel != "" {
			postMessageToUser(req.Text, userChannel, req.UserID)
		}
	}
}

// ExtractUsers parse message and return found users' ID
func getUsers(message string) []string {
	start, end := false, false
	return strings.FieldsFunc(message, func(r rune) bool {
		switch r {
		case '@':
			start, end = true, false
			return true
		case '|':
			start, end = false, true
		}
		if !start || end {
			return true
		}
		return false
	})
}
