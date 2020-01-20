package main

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"

	"github.com/tiqqe/go-logger"
)

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (res *events.APIGatewayProxyResponse, e error) {
	lctx, _ := lambdacontext.FromContext(ctx)
	logger.Init(lctx.AwsRequestID, os.Getenv("AWS_LAMBDA_FUNCTION_NAME"))

	defer handlePanic()

	logger.Info(&logger.LogEntry{
		Message: "Got event",
		Keys: map[string]interface{}{
			"Event": event,
		},
	})

	req := request{}
	res = &events.APIGatewayProxyResponse{StatusCode: 200}

	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {

		// ---------------------------------------
		// Command
		// ---------------------------------------
		logger.InfoString("Checking command")
		if result, ok := doCommand(event.Body); ok {
			logger.Info(&logger.LogEntry{
				Message: "Done Command",
				Keys: map[string]interface{}{
					"Body":   event.Body,
					"Result": result,
				},
			})
			res.Body = result
			return
		}

		// ---------------------------------------
		// Action
		// ---------------------------------------
		logger.InfoString("Checking action")
		if result, ok := doAction(event.Body); ok {
			logger.Info(&logger.LogEntry{
				Message: "Done Action",
				Keys: map[string]interface{}{
					"Body":   event.Body,
					"Result": result,
				},
			})
			res.Body = result
			return
		}

		logger.Error(&logger.LogEntry{
			Message:      "Failed to unmarshall request body",
			ErrorMessage: err.Error(),
			Keys:         map[string]interface{}{
				// "Event": event,
			},
		})

		res.Body = err.Error()
		return
	}

	logger.Info(&logger.LogEntry{
		Message: "Got request",
		Keys: map[string]interface{}{
			"ts":        req.Event.TS,
			"event_ts":  req.Event.EventTS,
			"delete_ts": req.Event.DeleteTS,
		},
	})

	if req.IsVerification() {
		res.Body = req.Challenge
		res.Headers = map[string]string{"Content-type": "text/plain"}
		logger.InfoString("Responding to challenge")
		return
	}

	if req.Event.Text == "" {
		res.Body = "Message is empty"
		logger.InfoString(res.Body)
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
		logger.InfoString(res.Body)
		return
	}

	threadChan := make(chan string)
	go getMainThread(req.Event.TS, req.Event.Channel, threadChan)

	go logger.InfoString("Translating message")

	body, err := translate(req.Event.Text)
	if err != nil {
		res.Body = err.Error()
		return
	}

	logger.Info(&logger.LogEntry{
		Message: "Got translation",
		Keys:    map[string]interface{}{
			// "Body": body,
		},
	})

	if strings.ToLower(body) == strings.ToLower(req.Event.Text) {
		logger.InfoString("Message is same as translation")
		return
	}

	logger.Info(&logger.LogEntry{
		Message: "Sending translated message to slack",
		Keys: map[string]interface{}{
			"Message":     req.Event.Text,
			"Translation": body,
		},
	})

	threadTs := <-threadChan

	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(user string) {
			defer wg.Done()
			err = postMessageToSlack(body, req.Event.Channel, req.Event.User, user, threadTs)
			if err != nil {
				logger.Error(&logger.LogEntry{
					Message:      "Failed to post slack message",
					ErrorMessage: err.Error(),
					Keys: map[string]interface{}{
						"User":    user,
						"Channel": req.Event.Channel,
					},
				})
				if err.Error() == "user_not_in_channel" {
					notInChannel[user+":"+req.Event.Channel] = true
				}
			}
		}(user)
	}
	wg.Wait()

	logger.InfoString("Done")

	return res, nil
}

func main() {
	lambda.Start(handler)
}

func handlePanic() {
	msg := recover()
	if msg != nil {
		entry := &logger.LogEntry{
			Message:   "Go panic",
			ErrorCode: "GoPanic",
		}
		switch msg := msg.(type) {
		case string:
			entry.ErrorMessage = msg
		case error:
			entry.ErrorMessage = msg.Error()

		default:
			entry.ErrorCode = "Unknown error type"
			entry.SetKey("error", msg)
		}

		logger.Error(entry)
	}
}
