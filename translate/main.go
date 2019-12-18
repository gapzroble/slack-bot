package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"

	"github.com/tiqqe/go-logger"
)

var translator = "https://translate.googleapis.com/translate_a/single?client=gtx&sl=sv&tl=en&dt=t&q="

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

	url := translator + url.QueryEscape(req.Event.Text)
	response, err := http.Get(url)
	if err != nil {
		logger.Error(&logger.LogEntry{
			Message:      "Failed to translate message from api",
			ErrorMessage: err.Error(),
			Keys: map[string]interface{}{
				"Url": url,
			},
		})
		res.Body = err.Error()
		return
	}

	if response.StatusCode > 299 {
		logger.Error(&logger.LogEntry{
			Message: "Expecting response status 2XX",
			Keys: map[string]interface{}{
				"Response": response,
			},
		})
		res.Body = "Translate response not OK"
		return
	}

	dat, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error(&logger.LogEntry{
			Message:      "Fail to read response body",
			ErrorMessage: err.Error(),
		})
		return
	}

	body := parse(string(dat))
	if body == "" {
		body = string(dat)
	}
	logger.InfoStringf("Got response: %s", body)

	if body == req.Event.Text {
		logger.InfoString("message is same as translation")
		return
	}

	logger.InfoStringf("Sending translation to slack: %s", body)

	var wg sync.WaitGroup

	for _, user := range users() {
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
					notInChannel[user+req.Event.Channel] = true
				}
			}
		}(user)
	}
	wg.Wait()

	logger.InfoString("Done")

	return res, nil
}

func users() []string {
	return strings.Split(slackUser, ",")
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
