package interpreter

import "testing"

func TestMessage_GetRegexEmpty(t *testing.T) {
	m := new(Message)
	re := m.GetRegex()
	if re.String() != "" {
		t.Error("Error in generating the regex. Expected:", "", "Got:", re.String())
	}
}

func TestMessage_GetRegex(t *testing.T) {
	m := new(Message)
	m.Formats = append(m.Formats, "hi")
	re := m.GetRegex()
	if re.String() != "(hi)" {
		t.Error("Error in generating the regex. Expected:", "(hi)", "Got:", re.String(), m)
	}
}

func TestMessage_GetRegexWithPostFix(t *testing.T) {
	m := new(Message)
	m.Formats = append(m.Formats, "hi")
	m.Postfixes = append(m.Postfixes, "!")
	re := m.GetRegex()
	if re.String() != "(hi)[!]{0,1}" {
		t.Error("Error in generating the regex. Expected:", "(hi)[!]{0,1}", "Got:", re.String())
	}
}

func TestMessage_GetRegexWithPreFix(t *testing.T) {
	m := new(Message)
	m.Formats = append(m.Formats, "hi")
	m.Prefixes = append(m.Prefixes, "!")
	re := m.GetRegex()
	if re.String() != "[!]{0,1}(hi)" {
		t.Error("Error in generating the regex. Expected:", "(hi)[!]{0,1}", "Got:", re.String())
	}
}
