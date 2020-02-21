package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (res *events.APIGatewayProxyResponse, e error) {
	defer handlePanic()

	log.Printf("Got event: %#v", event)

	req := request{}
	res = &events.APIGatewayProxyResponse{StatusCode: 200}

	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {

		// ---------------------------------------
		// Command
		// ---------------------------------------
		log.Print("Checking command")
		if result, ok := doCommand(event.Body); ok {
			log.Printf("Done command. Body: %s, Result: %s", event.Body, result)
			res.Body = result
			return
		}

		// ---------------------------------------
		// Action
		// ---------------------------------------
		log.Print("Checking action")
		if _, ok := doAction(event.Body); ok {
			log.Print("Done action.")
			// res.Body = result
			return
		}

		log.Printf("Failed to unmarshall request body: %s", err.Error())

		res.Body = err.Error()
		return
	}

	log.Printf("Got request. ts: %s, event_ts: %s, delete_ts: %s", req.Event.TS, req.Event.EventTS, req.Event.DeleteTS)

	if req.IsVerification() {
		res.Body = req.Challenge
		res.Headers = map[string]string{"Content-type": "text/plain"}
		log.Print("Responding to challenge")
		return
	}

	if req.Event.Text == "" {
		res.Body = "Message is empty"
		log.Print(res.Body)
		return
	}

	users := make([]string, 0, len(slackUsers))
	for user := range slackUsers {
		if user != req.Event.User && !userNotInChannel(user, req.Event.Channel) {
			users = append(users, user)
		}
	}
	if len(users) == 0 {
		res.Body = "No users to send to"
		log.Print(res.Body)
		return
	}

	threadTsChan := make(chan string, len(users))
	go func() {
		threadTs := getMainThread(req.Event.TS, req.Event.Channel)
		for range users {
			threadTsChan <- threadTs
		}
	}()

	senderPicChan := make(chan string, len(users))
	go func() {
		senderPic := getSenderPic(req.Event.User)
		for range users {
			senderPicChan <- senderPic
		}
	}()

	log.Print("Translating message")

	body, err := translate(req.Event.Text)
	if err != nil {
		res.Body = err.Error()
		return
	}

	log.Print("Got translation")

	if strings.ToLower(body) == strings.ToLower(req.Event.Text) {
		log.Print("Message is same as translation")
		return
	}

	log.Printf("Sending translated message to slack. Message: %s, Translation: %s", req.Event.Text, body)

	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(user string) {
			defer wg.Done()
			err = postMessageToSlack(body, req.Event.Channel, req.Event.User, user, threadTsChan, senderPicChan)
			if err != nil {
				log.Printf("Failed to post slack message: %s, User: %s, Channel: %s", err.Error(), user, req.Event.Channel)
				if err.Error() == "user_not_in_channel" {
					notInChannel[user+":"+req.Event.Channel] = true
				}
			}
		}(user)
	}
	wg.Wait()

	log.Print("Done")

	return res, nil
}

func main() {
	lambda.Start(handler)
}

func handlePanic() {
	msg := recover()
	if msg != nil {
		var err string
		switch msg := msg.(type) {
		case string:
			err = msg
		case error:
			err = msg.Error()

		default:
			err = fmt.Sprintf("Unknown error type: %#v", msg)
		}

		log.Printf("Go panic: %s", err)
	}
}
