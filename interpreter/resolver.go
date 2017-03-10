package interpreter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
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

func GetResponse(q string) (response string) {
	return ProcessQuery(q).Response
}

func ProcessQuery(q string) *Message {
	for k, v := range patternMap {
		fmt.Println(k, q, k.MatchString(q))
		if k.MatchString(q) {
			if v.Response != "" {
				return &v
			}
			switch v.Category {
			case "Show Justickets Order":
				res, err := GetOrder(q)
				if err != nil {
					log.Println("Error:", err)
					res = err.Error()
				}
				v.Response = res
				return &v
			default:

			}
		}
	}
	return &Message{}
}
