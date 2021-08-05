package gobot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type SlackEventHub struct {
	api          *slack.Client
	client       *socketmode.Client
	bangHandlers map[string]func(*MessageSource, string)
}

func NewSlackEventHub() (*SlackEventHub, error) {
	// 	appToken := os.Getenv("SLACK_APP_TOKEN")
	// 	if !strings.HasPrefix(appToken, "xapp-") {
	// 		return nil, errors.New("SLACK_APP_TOKEN must have the prefix \"xapp-\".")
	// 	}

	//     signingSecret := os.Getenv("SLACK_SIGNING_SECRET")
	//     if signingSecret == "" {
	//         return nil, errors.New("SLACK_SIGNING_SECRET must be set")
	//     }

	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if !strings.HasPrefix(botToken, "xoxb-") {
		return nil, errors.New("SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
	}

	// 	api := slack.New(
	// 		botToken,
	// 		// slack.OptionDebug(true),
	// 		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	// 		slack.OptionAppLevelToken(appToken),
	// 	)

	// 	client := socketmode.New(
	// 		api,
	// 		// socketmode.OptionDebug(true),
	// 		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	// 	)

	var api = slack.New(botToken)

	ret := &SlackEventHub{
		api: api,
		//client:       client,
		bangHandlers: make(map[string]func(*MessageSource, string)),
	}

	ret.RegisterBangHandler("hooks", ret.hooksHandler)
	return ret, nil
}

func (s *SlackEventHub) hooksHandler(source *MessageSource, msg string) {
	ret := make([]string, 0)
	for cmd, _ := range s.bangHandlers {
		if cmd != "hooks" {
			ret = append(ret, cmd)
		}
	}
	sort.Strings(ret)
	s.Message(source, strings.Join(ret, ", "))
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

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Got request: %v\n", request)

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
	} else {
		return events.APIGatewayProxyResponse{StatusCode: 500}, fmt.Errorf("Unrecognized event type: %v", eventsAPIEvent)
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func parseBody(body string) string {
	decodedValue, _ := url.QueryUnescape(body)
	data := strings.Trim(decodedValue, ":payload=")
	return data
}

func (s *SlackEventHub) StartEventLoop() {
	lambda.Start(Handler)
	// 	go func() {
	// 		for evt := range s.client.Events {
	// 			switch evt.Type {
	// 			case socketmode.EventTypeHello:
	// 				fmt.Println("Received hello")
	// 			case socketmode.EventTypeConnecting:
	// 				fmt.Println("Connecting to Slack with Socket Mode...")
	// 			case socketmode.EventTypeConnectionError:
	// 				fmt.Println("Connection failed. Retrying later...")
	// 			case socketmode.EventTypeConnected:
	// 				fmt.Println("Connected to Slack with Socket Mode.")
	// 			case socketmode.EventTypeEventsAPI:
	// 				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	// 				if !ok {
	// 					fmt.Printf("Event not ok %+v\n", evt)
	// 					continue
	// 				}

	// 				// fmt.Printf("Event received: %+v\n", eventsAPIEvent)

	// 				s.client.Ack(*evt.Request)

	// 				switch eventsAPIEvent.Type {
	// 				case slackevents.CallbackEvent:
	// 					innerEvent := eventsAPIEvent.InnerEvent
	// 					switch ev := innerEvent.Data.(type) {
	// 					case *slackevents.MessageEvent:
	// 						s.handleBangs(ev)
	// 						// fmt.Printf("Message received: %+v", ev)
	// 					default:
	// 						fmt.Printf("Unsupported inner event type %T\n", innerEvent.Data)
	// 					}
	// 				default:
	// 					fmt.Printf("Unsupported Events API event type %s\n", eventsAPIEvent.Type)
	// 				}
	// 			default:
	// 				fmt.Fprintf(os.Stderr, "Unsupported event type received: %s\n", evt.Type)
	// 			}
	// 		}
	// 	}()

	// 	s.client.Run()
}
