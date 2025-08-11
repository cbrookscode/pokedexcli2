package repl

import (
	"strings"
	"unicode"
)

func CleanInput(text string) []string {
	cleaned := []string{}
	for _, word := range strings.Fields(strings.ToLower(text)) {
		newword := ""
		for _, letter := range word {
			if unicode.IsLetter(letter) {
				newword = newword + string(letter)
			}
		}
		cleaned = append(cleaned, newword)
	}
	return cleaned
}
