package interpreter

import (
	"regexp"
	"strings"
)

type Message struct {
	ID        string   `json:"id"`
	Category  string   `json:"category"`
	Prefixes  []string `json:"prefixes"`
	Case      bool     `json:"case"`
	Formats   []string `json:"formats"`
	Postfixes []string `json:"postfixes"`
	Response  string   `json:"response"`
	Regex     *regexp.Regexp
}

func (m *Message) GetRegex() *regexp.Regexp {
	if m.Regex != nil {
		return m.Regex
	}

	re := ""
	if m.Case {
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
		re += "["
		for _, v := range m.Prefixes {
			re += v + "|"
		}
		re = strings.TrimRight(re, "|")
		re += "]{0,1}"
	}
	return re
}

func (m *Message) getPostfixRegex() string {
	re := ""
	if len(m.Postfixes) > 0 {
		re += "["
		for _, v := range m.Postfixes {
			re += v + "|"
		}
		re = strings.TrimRight(re, "|")
		re += "]{0,1}"
	}
	return re
}
