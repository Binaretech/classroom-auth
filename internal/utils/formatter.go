package utils

import "strings"

func LowerCaseInitial(source string) string {
	return strings.ToLower(string(source[0])) + string(source[1:])
}
