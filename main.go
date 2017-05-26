package main

import (
	"fmt"
	"log"
	"net/http"
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
	api.SetDebug(false)

	rtm := api.NewRTM()
	go rtm.ManageConnection()
	go runSlackListener(rtm, api)
	router := GetRouter()
	log.Println(os.Getenv("APP_NAME"), "listening in port", os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), router)
}

func respond(rtm *slack.RTM, api *slack.Client, msg *slack.MessageEvent, prefix string) {
	text := msg.Text
	text = strings.Replace(text, prefix, "", -1)
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	response := interpreter.ProcessQuery(text)
	if response.Attachments[0].Pretext == "" {
		rtm.SendMessage(rtm.NewOutgoingMessage(
			`I'm sorry, I don't understand! Sometimes I have an easier time with a few simple keywords.`,
			msg.Channel))
	} else {
		api.PostMessage(msg.Channel, "", response)
	}
}

func runSlackListener(rtm *slack.RTM, api *slack.Client) {

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Println("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				info := rtm.GetInfo()
				prefix := fmt.Sprintf("<@%s>", info.User.ID)
				//TODO add prefix exception if message is a direct message to bot
				if ev.User != info.User.ID && strings.Contains(ev.Text, prefix) {
					go respond(rtm, api, ev, prefix)
				}
				//fmt.Println(ev.User != info.User.ID && strings.HasPrefix(ev.Text, prefix))
			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:

			}
		}
	}
}
