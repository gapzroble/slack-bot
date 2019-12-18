package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

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

	dat, _ := ioutil.ReadAll(response.Body)
	body := string(dat)
	logger.InfoStringf("Got response: %s", body)

	body = strings.Split(body, `","`+req.Event.Text)[0]
	parts := strings.Split(body, `["`)
	if len(parts) <= 1 {
		logger.Error(&logger.LogEntry{
			Message: "Error parsing response",
			Keys: map[string]interface{}{
				"Parts": parts,
			},
		})
		return
	}
	body = parts[1]
	if body == req.Event.Text {
		logger.InfoString("message is same as translation")
		return
	}

	user := slackUser
	if len(req.AuthedUsers) > 0 {
		user = req.AuthedUsers[0]
	}

	if user == req.Event.User {
		logger.InfoString("You are the sender")
		return
	}

	logger.InfoStringf("Sending translation to slack: %s", body)
	err = postMessageToSlack(body, req.Event.Channel, req.Event.User, user, req.Event.EventTS)
	if err != nil {
		logger.Error(&logger.LogEntry{
			Message:      "Failed to post slack message",
			ErrorMessage: err.Error(),
			Keys: map[string]interface{}{
				"Message": body,
			},
		})
		res.Body = err.Error()
		return res, nil
	}
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
