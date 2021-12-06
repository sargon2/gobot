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

type LambdaEventProcessor struct {
	api             *slack.Client
	messageCallback func(*slackevents.MessageEvent) // TODO this should be a list of callbacks
}

func NewLambdaEventProcessor() *LambdaEventProcessor {
	return &LambdaEventProcessor{}
}

func (s *LambdaEventProcessor) GetUsername(userID string) string {
	user, err := s.api.GetUserInfo(userID)
	if err != nil {
		fmt.Println(err.Error())
	}
	if err != nil || user == nil {
		return userID
	}
	return user.Name
}

func (s *LambdaEventProcessor) RegisterMessageCallback(cb func(*slackevents.MessageEvent)) {
	s.messageCallback = cb
}

func (s *LambdaEventProcessor) Message(source *MessageSource, m string) {
	fmt.Println("Message", m)
	_, _, err := s.api.PostMessage(source.ChannelID, slack.MsgOptionText(m, true))
	if err != nil {
		fmt.Println(err)
	}
}

func (s *LambdaEventProcessor) HandleEvent(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// fmt.Printf("Got request: %v\n", request)
	if _, ok := request.Headers["X-Slack-Retry-Num"]; ok {
		fmt.Println("X-Slack-Retry-Num set, aborting")
		return events.APIGatewayProxyResponse{StatusCode: 200}, nil
	}

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

func (s *LambdaEventProcessor) StartProcessingEvents() error {
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if !strings.HasPrefix(botToken, "xoxb-") {
		return errors.New("SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
	}

	s.api = slack.New(botToken)

	lambda.Start(s.HandleEvent)
	return nil
}

func (*LambdaEventProcessor) IsTestMode() bool {
	return false
}
