package interpreter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"strconv"

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
				res, err := GetOrder(q)
				if err != nil {
					log.Println("Error:", err)
					attachment.Pretext = err.Error()
					return params
				}
				if res.SessionID == "" {
					attachment.Pretext = "No Order found for order id " + res.BlockCode
					attachment.Color = "#BDB76B"
					attachment.Fields = append(attachment.Fields,
						slack.AttachmentField{
							Title: "Seach On Admin",
							Value: "https://admin.justickets.co/bookings?detail=" + res.BlockCode,
							Short: true})
					return params
				}
				attachment.Pretext = "Found Order: " + res.BlockCode
				attachment.Fields = append(attachment.Fields, slack.AttachmentField{
					Title: "Name",
					Value: res.Name,
					Short: true})
				attachment.Fields = append(attachment.Fields, slack.AttachmentField{Title: "Email",
					Value: res.Email,
					Short: true})
				attachment.Fields = append(attachment.Fields, slack.AttachmentField{Title: "Mobile",
					Value: res.Mobile,
					Short: true})
				attachment.Fields = append(attachment.Fields, slack.AttachmentField{Title: "Paid",
					Value: strconv.FormatBool(res.Paid),
					Short: true})
				attachment.Fields = append(attachment.Fields, slack.AttachmentField{Title: "Confirmed",
					Value: strconv.FormatBool(res.Confirmed),
					Short: true})
				attachment.Fields = append(attachment.Fields, slack.AttachmentField{Title: "Failed",
					Value: strconv.FormatBool(res.Cancelled),
					Short: true})
				attachment.Fields = append(attachment.Fields, slack.AttachmentField{Title: "Bill Total",
					Value: strconv.FormatFloat(res.Bill.Total, 'G', 6, 64),
					Short: true})
				return params
			default:

			}
		}
	}
	return params
}
