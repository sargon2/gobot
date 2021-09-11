package gobot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type SlackEventHub struct {
	api             *slack.Client
	messageCallback func(*slackevents.MessageEvent)
}

func NewSlackEventHub() (*SlackEventHub, error) {
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if !strings.HasPrefix(botToken, "xoxb-") {
		return nil, errors.New("SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
	}

	var api = slack.New(botToken)

	ret := &SlackEventHub{
		api: api,
	}

	return ret, nil
}

func (s *SlackEventHub) RegisterMessageCallback(cb func(*slackevents.MessageEvent)) {
	s.messageCallback = cb
}

func (s *SlackEventHub) Message(source *MessageSource, m string) {
	_, _, err := s.api.PostMessage(source.ChannelID, slack.MsgOptionText(m, true))
	if err != nil {
		fmt.Println(err) // TODO proper logging?
	}
}

func (s *SlackEventHub) HandleEvent(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// fmt.Printf("Got request: %v\n", request)

	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(request.Body), slackevents.OptionNoVerifyToken())
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(request.Body), &r)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500}, err
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "text"},
			Body:       r.Challenge,
		}, nil
	} else if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			// fmt.Printf("Message received: %+v", ev)
			s.messageCallback(ev)
		default:
			fmt.Printf("Unsupported inner event type %T\n", innerEvent.Data)
		}
	} else {
		return events.APIGatewayProxyResponse{StatusCode: 500}, fmt.Errorf("Unrecognized event type: %v", eventsAPIEvent)
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func (s *SlackEventHub) StartEventLoop() {
	lambda.Start(s.HandleEvent)
}
