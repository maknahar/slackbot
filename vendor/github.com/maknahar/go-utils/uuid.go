package go_utils

import "regexp"

var RegexUUIDV1 = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-1[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
var RegexUUIDV2 = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-2[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
var RegexUUIDV3 = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-3[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
var RegexUUIDV4 = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
var RegexUUIDV5 = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-5[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

func IsValidUUIDV1(uuid string) bool {
	return RegexUUIDV1.MatchString(uuid)
}

func IsValidUUIDV2(uuid string) bool {
	return RegexUUIDV2.MatchString(uuid)
}

func IsValidUUIDV3(uuid string) bool {
	return RegexUUIDV3.MatchString(uuid)
}

func IsValidUUIDV4(uuid string) bool {
	return RegexUUIDV4.MatchString(uuid)
}

func IsValidUUIDV5(uuid string) bool {
	return RegexUUIDV5.MatchString(uuid)
}
