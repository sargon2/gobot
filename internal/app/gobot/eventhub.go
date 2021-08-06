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
	api          *slack.Client
	bangHandlers map[string]func(*MessageSource, string)
}

func NewSlackEventHub() (*SlackEventHub, error) {
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if !strings.HasPrefix(botToken, "xoxb-") {
		return nil, errors.New("SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
	}

	var api = slack.New(botToken)

	ret := &SlackEventHub{
		api: api,
		bangHandlers: make(map[string]func(*MessageSource, string)),
	}

	return ret, nil
}

func (s *SlackEventHub) Message(source *MessageSource, m string) {
	_, _, err := s.api.PostMessage(source.ChannelID, slack.MsgOptionText(m, true))
	if err != nil {
		fmt.Println(err) // TODO proper logging?
	}
}

func (s *SlackEventHub) RegisterBangHandler(cmd string, handler func(*MessageSource, string)) {
	s.bangHandlers[cmd] = handler
}

func (s *SlackEventHub) GetBangHandlers() (map[string]func(*MessageSource, string)) {
    return s.bangHandlers
}

func (s *SlackEventHub) handleBangs(event *slackevents.MessageEvent) {
	messageText := event.Text
	channelID := event.Channel
	for cmd, handler := range s.bangHandlers {
		bangCmd := "!" + cmd
		if messageText == bangCmd || strings.HasPrefix(messageText, bangCmd+" ") {
			source := &MessageSource{
				ChannelID: channelID,
			}
			messageText = strings.TrimSpace(messageText[len(bangCmd):])
			handler(source, messageText)
		}

	}
}

func (s *SlackEventHub) Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) { // TODO rename function
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
			s.handleBangs(ev)
		default:
			fmt.Printf("Unsupported inner event type %T\n", innerEvent.Data)
		}
	} else {
		return events.APIGatewayProxyResponse{StatusCode: 500}, fmt.Errorf("Unrecognized event type: %v", eventsAPIEvent)
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func (s *SlackEventHub) StartEventLoop() {
	lambda.Start(s.Handler)
}
