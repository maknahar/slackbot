package interpreter

import (
	"regexp"
	"strings"

	"github.com/nlopes/slack"
)

type Message struct {
	ID               string   `json:"id"`
	Category         string   `json:"category"`
	Prefixes         []string `json:"prefixes"`
	PrefixMendatory  bool     `json:"prefixMandatory"`
	Case             bool     `json:"case"`
	Formats          []string `json:"formats"`
	Postfixes        []string `json:"postfixes"`
	PostfixMendatory bool     `json:"postfixMandatory"`
	Response         string   `json:"response"`
	Regex            *regexp.Regexp
}

func GetSlackMessage() slack.PostMessageParameters {
	msg := slack.PostMessageParameters{}
	msg.AsUser = true
	msg.Attachments = append(msg.Attachments, slack.Attachment{
		Color:      "#9932CC",
		AuthorName: "Justickets Bot",
		//AuthorSubname: "Mayank Patel",
		AuthorLink: "https://github.com/maknahar/jtbot",
		AuthorIcon: "https://data.justickets.co/favicon.ico",
		Footer:     "Always at your service",
		FooterIcon: "http://cconnect.s3.amazonaws.com/wp-content/uploads/2017/02/2017-Funko-Pop-Mystery-Science-Theater-3000-Crow-T-Robot-e1486480774184.jpg",
	})
	return msg
}

func (m *Message) GetRegex() *regexp.Regexp {
	if m.Regex != nil {
		return m.Regex
	}

	re := ""
	if !m.Case {
		re += "(?i)"
	}
	re += m.getPrefixRegex()
	if len(m.Formats) > 0 {
		re += "("
		for _, v := range m.Formats {
			re += v + "|"
		}
		re = strings.TrimRight(re, "|")
		re += ")"
	}
	re += m.getPostfixRegex()
	m.Regex = regexp.MustCompile(re)

	return m.Regex
}

func (m *Message) getPrefixRegex() string {
	re := ""
	if len(m.Prefixes) > 0 {
		re += "("
		for _, v := range m.Prefixes {
			re += v + "|"
		}
		re = strings.TrimRight(re, "|")
		re += ")"
		if m.PrefixMendatory {
			re += "{1,1}"
		} else {
			re += "{0,1}"
		}
	}
	return re
}

func (m *Message) getPostfixRegex() string {
	re := ""
	if len(m.Postfixes) > 0 {
		re += "("
		for _, v := range m.Postfixes {
			re += v + "|"
		}
		re = strings.TrimRight(re, "|")
		re += ")"
		if m.PostfixMendatory {
			re += "{1,1}"
		} else {
			re += "{0,1}"
		}
	}
	return re
}
