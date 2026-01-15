package cleaninput

import "strings"

func CleanInput(text string) []string {
	elements := strings.Fields(strings.ToLower(text))
	return elements
}
