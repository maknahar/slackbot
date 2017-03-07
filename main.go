package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/maknahar/jtbot/interpreter"
	"github.com/nlopes/slack"
)

var (
	SLACK_TOKEN = os.Getenv("SLACK_TOKEN")
)

func main() {

	if SLACK_TOKEN == "" {
		log.Fatal("Missing SLACK_TOKEN env var")
	}

	api := slack.New(SLACK_TOKEN)
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Println("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", ev)
				info := rtm.GetInfo()
				prefix := fmt.Sprintf("<@%s> ", info.User.ID)

				if ev.User != info.User.ID && strings.HasPrefix(ev.Text, prefix) {
					respond(rtm, ev, prefix)
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
				//Take no action
			}
		}
	}
}

func respond(rtm *slack.RTM, msg *slack.MessageEvent, prefix string) {
	var response string
	text := msg.Text
	text = strings.TrimPrefix(text, prefix)
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	response = interpreter.GetResponse(text)
	if response == "" {
		rtm.SendMessage(rtm.NewOutgoingMessage(`I'm sorry, I don't understand! Sometimes I have an easier time with a few simple keywords.`, msg.Channel))
	} else {
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
	}
}
