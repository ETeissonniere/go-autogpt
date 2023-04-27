package helpers

import (
	"regexp"
	"strings"
)

func SplitCommand(input string) []string {
	re := regexp.MustCompile(`[^\s"']+|"[^"]*"|'[^']*'`)
	result := re.FindAllString(input, -1)
	for i, s := range result {
		result[i] = strings.Trim(s, "'\"")
	}

	return result
}

func WithEscapeCharacters(s string) string {
	escapeSequences := map[string]string{
		"\\a":  "\a",
		"\\b":  "\b",
		"\\f":  "\f",
		"\\n":  "\n",
		"\\r":  "\r",
		"\\t":  "\t",
		"\\v":  "\v",
		"\\\\": "\\",
		"\\'":  "'",
		"\\\"": "\"",
	}

	for escape, replacement := range escapeSequences {
		s = strings.ReplaceAll(s, escape, replacement)
	}

	return s
}
