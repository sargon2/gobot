package gobot

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type SlackSocketHub struct {
	api             *slack.Client
	client          *socketmode.Client
	messageCallback func(*slackevents.MessageEvent)
}

func NewSlackSocketHub() (*SlackSocketHub, error) {
	appToken := os.Getenv("SLACK_APP_TOKEN")
	if !strings.HasPrefix(appToken, "xapp-") {
		return nil, errors.New("SLACK_APP_TOKEN must have the prefix \"xapp-\".")
	}

	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if !strings.HasPrefix(botToken, "xoxb-") {
		return nil, errors.New("SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
	}

	api := slack.New(
		botToken,
		// slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(appToken),
	)

	client := socketmode.New(
		api,
		// socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	ret := &SlackSocketHub{
		api:    api,
		client: client,
	}

	return ret, nil
}

func (s *SlackSocketHub) RegisterMessageCallback(cb func(*slackevents.MessageEvent)) {
	s.messageCallback = cb
}

func (s *SlackSocketHub) Message(source *MessageSource, m string) {
	_, _, err := s.api.PostMessage(source.ChannelID, slack.MsgOptionText(m, true))
	if err != nil {
		fmt.Println(err)
	}
}

func (s *SlackSocketHub) StartEventLoop() {
	go func() {
		for evt := range s.client.Events {
			switch evt.Type {
			case socketmode.EventTypeHello:
				fmt.Println("Received hello")
			case socketmode.EventTypeConnecting:
				fmt.Println("Connecting to Slack with Socket Mode...")
			case socketmode.EventTypeConnectionError:
				fmt.Println("Connection failed. Retrying later...")
			case socketmode.EventTypeConnected:
				fmt.Println("Connected to Slack with Socket Mode.")
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					fmt.Printf("Event not ok %+v\n", evt)
					continue
				}

				// fmt.Printf("Event received: %+v\n", eventsAPIEvent)

				s.client.Ack(*evt.Request)

				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					innerEvent := eventsAPIEvent.InnerEvent
					switch ev := innerEvent.Data.(type) {
					case *slackevents.MessageEvent:
						s.messageCallback(ev)
						// fmt.Printf("Message received: %+v", ev)
					default:
						fmt.Printf("Unsupported inner event type %T\n", innerEvent.Data)
					}
				default:
					fmt.Printf("Unsupported Events API event type %s\n", eventsAPIEvent.Type)
				}
			default:
				fmt.Fprintf(os.Stderr, "Unsupported event type received: %s\n", evt.Type)
			}
		}
	}()

	s.client.Run()
}
