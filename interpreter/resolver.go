package interpreter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/nlopes/slack"
)

var patternMap map[*regexp.Regexp]Message

func init() {
	databaseFile := os.Getenv("DB_FILE")
	if databaseFile != "" {
		d, err := ioutil.ReadFile(databaseFile)
		if err != nil {
			log.Println("Error in reading database file")
			panic(err)
		}
		var databaseFile struct {
			Messages []Message `json:"messages"`
		}

		err = json.Unmarshal(d, &databaseFile)
		if err != nil {
			log.Println("Error in decoding database file")
			panic(err)
		}
		patternMap = make(map[*regexp.Regexp]Message)
		for _, v := range databaseFile.Messages {
			patternMap[v.GetRegex()] = v
		}
	}
}

func ProcessQuery(q string) slack.PostMessageParameters {
	params := GetSlackMessage()
	attachment := &params.Attachments[0]

	for k, v := range patternMap {
		fmt.Println(k, q, k.MatchString(q))
		if k.MatchString(q) {
			if v.Response != "" {
				attachment.Pretext = v.Response
				return params
			}
			switch v.Category {
			case "Show Justickets Order":
				order, err := GetOrder(q)
				if err != nil {
					log.Println("Error:", err)
					attachment.Pretext = err.Error()
					return params
				}
				order.FormatSlackMessage(attachment)
				return params
			case "Show Justickets Bill":
				order, err := GetOrder(q)
				if err != nil {
					log.Println("Error:", err)
					attachment.Pretext = err.Error()
					return params
				}
				order.FormatSlackMessageForBill(attachment)
				return params

			case "Staging Report is Down":
				r, err := GetReportStatus(true)
				if err != nil {
					log.Println("Error:", err)
					attachment.Pretext = err.Error()
					return params
				}
				if r == nil {
					attachment.Pretext = "Report is in sync"
					return params
				}
				cause := r.GetDelayReason()
				if cause == "" {
					cause = "Sorry, I could not diagnose the problem in report sync delay"
				}
				attachment.Pretext = cause

				return params

			case "Report is Down":
				r, err := GetReportStatus(false)
				if err != nil {
					log.Println("Error:", err)
					attachment.Pretext = err.Error()
					return params
				}
				if r == nil {
					attachment.Pretext = "Report is in sync"
					return params
				}
				cause := r.GetDelayReason()
				if cause == "" {
					cause = "Sorry, I could not diagnose the problem in report sync delay"
				}
				attachment.Pretext = cause

				return params
			default:

			}
		}
	}
	return params
}

func FormatSlackMessageReport(attachment *slack.Attachment) {
	attachment.Pretext = "Staging report is down"
}
