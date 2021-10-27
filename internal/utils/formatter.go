package utils

import "strings"

// LowerCaseInitial returns the given string with the first character in lower case.
func LowerCaseInitial(source string) string {
	return strings.ToLower(string(source[0])) + string(source[1:])
}
