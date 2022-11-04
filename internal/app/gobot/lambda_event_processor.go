package gobot

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type LambdaEventProcessor struct {
	api             *slack.Client
	messageCallback func(*slackevents.MessageEvent) // TODO this should be a list of callbacks
}

func NewLambdaEventProcessor(slackClient *slack.Client) (*LambdaEventProcessor, error) {
	return &LambdaEventProcessor{api: slackClient}, nil
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
	fmt.Printf("Got request: %+v\n", request)

	// Slack will retry if we don't respond in 3 seconds.
	// But we do respond after that, so we don't want retries.
	// So as a workaround we can ignore retried requests.
	// TODO The real fix for this would be to return 200 immediately, and use step functions to execute the bot after that.
	// https://stackoverflow.com/a/44670387
	if _, ok := request.Headers["x-slack-retry-num"]; ok {
		fmt.Println("x-slack-retry-num set, aborting")
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
	lambda.Start(s.HandleEvent)
	return nil
}

func (*LambdaEventProcessor) IsTestMode() bool {
	return false
}
