package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func hanleWithout(c *gin.Context) {
	var req with

	if err := c.ShouldBindWith(&req, binding.FormPost); err != nil {
		fmt.Println("Bind error:", err)
		c.JSON(500, err)
		return
	}

	fmt.Printf("req: %+v\n", req)

	c.Data(200, "text/plain", nil)

	users := getUsers(req.Text)
	if len(users) == 0 {
		return
	}

	for _, user := range users {
		if userChannel := gc.removeFromGc(req.ChannelID, user); userChannel != "" {
			postMessageToUser(req.Text, userChannel, req.UserID)
		}
	}
}
