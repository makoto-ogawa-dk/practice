package api

import "regexp"

var monthPattern = regexp.MustCompile(`^[0-9]{4}-(0[1-9]|1[0-2])$`)

func isValidMonth(v string) bool {
	return monthPattern.MatchString(v)
}
