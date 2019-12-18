package main

import (
	"context"
	"encoding/json"
	"os"
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

		// check here to re-use global var
		if result, ok := command(event.Body); ok {
			logger.InfoString("Done command")
			res.Body = result
			return
		}

		logger.Error(&logger.LogEntry{
			Message:      "Failed to unmarshall request body",
			ErrorMessage: err.Error(),
			Keys: map[string]interface{}{
				"Event": event,
			},
		})
		res.Body = err.Error()
		return
	}

	logger.Info(&logger.LogEntry{
		Message: "Got request",
		Keys: map[string]interface{}{
			"Request": req,
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

	logger.InfoString("Translating message")

	body, err := translate(req.Event.Text)
	if err != nil {
		res.Body = err.Error()
		return
	}

	logger.InfoStringf("Got response: %s", body)

	if body == req.Event.Text {
		logger.InfoString("message is same as translation")
		return
	}

	logger.InfoStringf("Sending translation to slack: %s", body)

	var wg sync.WaitGroup

	for user := range slackUsers {
		if user == req.Event.User || userNotInChannel(user, req.Event.Channel) {
			continue
		}
		wg.Add(1)
		go func(user string) {
			defer wg.Done()
			err = postMessageToSlack(body, req.Event.Channel, req.Event.User, user, req.Event.EventTS)
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
