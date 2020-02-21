package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	team "github.com/rroble/slack-bot/groupChat/users"
	"github.com/rroble/slack-bot/lib/slack"
)

var (
	slackToken string
	botToken   string
	users      team.Team
	gc         groupChat
)

func init() {
	slackToken = os.Getenv("SLACK_TOKEN")
	botToken = os.Getenv("BOT_TOKEN")
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (res *events.APIGatewayProxyResponse, e error) {
	defer handlePanic()

	log.Printf("Got event: %#v", event)

	req := slack.Request{}
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

	threadTsChan := make(chan string, len(users))
	go func() {
		threadTs := slack.GetMainThread(req.Event.TS, req.Event.Channel)
		for range users {
			threadTsChan <- threadTs
		}
	}()

	senderPicChan := make(chan string, len(users))
	go func() {
		senderPic := slack.GetSenderPic(req.Event.User)
		for range users {
			senderPicChan <- senderPic
		}
	}()

	go func() {
		users, _ = team.Load(botToken)
		fmt.Println("users loaded:", len(users))
	}()

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
