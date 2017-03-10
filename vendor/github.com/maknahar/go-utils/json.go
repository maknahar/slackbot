package go_utils

import (
	"bytes"
	"encoding/json"
)

//JsonPrettyPrint returns the given json in pretty format.
// If it fails to pretty print the JSON it just returns the original string.
// Useful for printing HTTP responses that should contain JSON.
// If indent is given as empty tabs are used
func JsonPrettyPrint(in string, prefix, indent string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}
